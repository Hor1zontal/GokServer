package object

import (
	"fmt"
	"net"
	"os"
	"stresstest/base"
	"aliens/common/cipher/xxtea"
	"encoding/binary"
	"github.com/gogo/protobuf/proto"
	"sync"
	"gok/service/msg/protocol"
	"time"
	"stresstest/testcase"
	"aliens/log"
	"stresstest/message"
	"bytes"
)

var TEST_COUNT int = 0

const (
	RET_TYPE_DONE	int32 = 1
	RET_TYPE_PUSH	int32 = 2
	RET_TYPE_RESULT	int32 = 3
	RET_TYPE_ERROR	int32 = 4
)

type Player struct {
	sync.RWMutex
	m_Idx int32 //索引

	syncStartTime time.Time

	seq      int32
	Uid      int32
	Token    string
	Username string
	Password string
	accountType	int32
	accountNum int

	Game      bool
	synctime  int64
	secretkey string

	m_SrvInfo string //游戏服务器信息
	passport  string //
	game      string

	m_pCon *net.TCPConn //tcp 的连接器

	m_bReadyDisCon bool //是否准备断开连接

	m_Channel        chan interface{}
	m_Channel_Closed bool //通道是否 已经关闭了

	m_State int32 //玩家状态

	testSuites []testcase.TestSuite
	currTestSuiteIndex int

	starID			int32
	starType		int32
}
//用户数据
//type User struct {
//	userID		int32
//	level		int32
//	exp			int32
//	power		int32
//	faith 		int32
//	items		[]*Item
//	currentStar	*Star
//	stars		[]*Star
//}
//
//type Item struct {
//
//}
//
//type Star struct {
//	starID 		int32
//	starType	int32
//	ownID		int32
//	buildings	[]*Building
//	itemGroups 	[]*ItemGroup
//}
//
//type Building struct {
//	buildID			int32
//	buildType 		int32
//	buildLevel		int32
//}
//
//type ItemGroup struct {
//	groupID			int32
//	done 			bool
//	active			bool
//}

/**
*	@brief 获取角色索引
 */
func (this *Player) GetIdx() int32 { return this.m_Idx }

/**
*	@brief 初始化
 */
func (this *Player) Init(idx int32, passportServer string, gameServer string, synctime int64, secretkey string, accountType int32, accountTotalNum int) bool {
	this.m_Idx = idx
	this.synctime = synctime
	this.m_SrvInfo = passportServer
	this.passport = passportServer
	this.secretkey = secretkey
	this.game = gameServer
	this.accountType = accountType
	this.accountNum = accountTotalNum
	this.testSuites = []testcase.TestSuite{}
	this.m_State = base.PLAYER_STATE_NONE //默认设置为 无状态
	this._open_channel()                  //打开通道
	if this.IsConnect() == false {
		//fmt.Printf(">>> try to connect \n");
		if this.Connect() == false { //连接失败
			fmt.Printf(">>> connect failed %v\n", this.m_Idx)
			return false
		}
		fmt.Printf(">>> connect success %v\n", this.m_Idx)
	}


	return true
}

func (this *Player) AcceptResult(result *protocol.GS2C) {
	if this.currTestSuiteIndex >= len(this.testSuites) {
		return
	}
	testSuite := this.testSuites[this.currTestSuiteIndex]
	testSuite.AcceptResult(result)
}

func (this *Player) NextMessage() *protocol.C2GS {
	if this.currTestSuiteIndex >= len(this.testSuites) {
		return nil
	}
	testSuite := this.testSuites[this.currTestSuiteIndex]
	message := testSuite.NextMessage()
	if message != nil {
		return message
	}
	this.currTestSuiteIndex++
	return this.NextMessage()
}



func (this *Player) ReconnectGame() bool {
	//已经连接了这个地址不需要重复连接
	if this.IsConnect() && this.m_SrvInfo == this.game {
		return true
	}
	this.m_SrvInfo = this.game
	this.m_State = base.PLAYER_STATE_NONE //默认设置为 无状态
	this._open_channel()                  //打开通道

	if this.IsConnect() {
		this.DisConnect()
		fmt.Printf(">>> reconnect success %v \n", this.m_Idx)
	}
	fmt.Printf(">>> try to connect %v\n", this.m_Idx)
	if this.Connect() == false { //连接失败
		fmt.Printf(">>> reconnect failed %v\n", this.m_Idx)
		return false
	} else {
		fmt.Printf(">>> reconnect success %v\n", this.m_Idx)
	}
	return true
}

func (this *Player) ReconnectPassport() bool {
	//已经连接了这个地址不需要重复连接
	if this.IsConnect() && this.m_SrvInfo == this.passport {
		return true
	}
	this.m_SrvInfo = this.passport
	this.m_State = base.PLAYER_STATE_NONE //默认设置为 无状态
	this._open_channel()                  //打开通道

	if this.IsConnect() {
		this.DisConnect()
		fmt.Printf(">>> reconnect success %v \n", this.m_Idx)
	}
	fmt.Printf(">>> try to connect %v\n", this.m_Idx)
	if this.Connect() == false { //连接失败
		fmt.Printf(">>> reconnect failed %v\n", this.m_Idx)
		return false
	} else {
		fmt.Printf(">>> reconnect success %v\n", this.m_Idx)
	}
	return true
}

/**
*	@brief 插入操作码
 */
func (this *Player) AcceptOp(op int) {
	if this.m_Channel == nil {
		return
	}
	select {
	case this.m_Channel <- op:
	default:
		fmt.Printf("%d message channel full\n", this.m_Idx)
		//TODO 消息管道满了需要通知客户端消息请求太过频繁
	}
}

func (this *Player) isCipher() bool {
	return this.secretkey != ""
}

/**
*	@brief 请求连接
 */
func (this *Player) Connect() bool {
	this.Lock()
	defer this.Unlock()
	if this.m_SrvInfo == "" {
		return false
	}
	if this.m_pCon != nil {
		return false
	}
	addr, err := net.ResolveTCPAddr("tcp", this.m_SrvInfo)
	if err != nil {
		fmt.Printf("Error: Idx=%d, %s\n", this.m_Idx, err.Error())
		return false
	}
	this.m_pCon, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Printf("Error: Idx=%d, %s\n", this.m_Idx, err.Error())
		this.m_pCon = nil //重新设置为nil
		return false
	}

	buf := make([]byte, 8182)
	go func() {
		for {
			if this.m_pCon == nil {
				break
			}
			bufLen, err := this.m_pCon.Read(buf)
			if err != nil {
				this.Lock()
				defer this.Unlock()
				this.m_pCon = nil
				fmt.Printf("read error: %v\n", err.Error())
				break
			}
			receiveBuf := buf[:bufLen]
			//fmt.Printf("Recv msg = %v\n   %v",  buf )
			thisData, nextData := this.splitMessage(receiveBuf)

			for {
				_, ret := this.acceptMessageData(thisData)
				if ret == RET_TYPE_DONE {
					TEST_COUNT++
					if TEST_COUNT == this.accountNum {
						log.Fatal("testSuite execute times:%v",TEST_COUNT)
					}
				}
				if ret == RET_TYPE_RESULT {
					log.Fatal("###########Test Result Push")
				}
				if len(nextData) == 0 {
					break
				}
				thisData , nextData = this.splitMessage(nextData)
			}
			//if code == 11024{
			//	this.m_State = base.PLAYER_STATE_GAME;//表示 玩家已经进入游戏
			//}
		}
	}()

	return true
}

func (this *Player) splitMessage(data []byte) ([]byte, []byte){
	dataLen := BytesToInt(data[:2])
	thisData := data[2:dataLen+2]
	nextData := data[len(thisData)+2:]
	return thisData, nextData
}

func  (this *Player) acceptMessageData(data []byte) (*protocol.GS2C, int32) {
	recv := &protocol.GS2C{}
	if this.isCipher() {
		data = xxtea.Decrypt(data, []byte(this.secretkey))
	}
	//log.Debug("receive marshal message %v", data)
	err := proto.Unmarshal(data, recv)
	if err != nil {
		fmt.Printf("unmarshaling error: %v\n", err.Error())
		return nil, RET_TYPE_ERROR
	}
	if recv.Sequence[0] > 1000 {
		fmt.Printf("%v => push   : %v %v\n", this.Username, recv, time.Now())
		return nil, RET_TYPE_PUSH
	}
	ret := this._receive_Message(recv)
	return recv, ret
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int16
	binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
	return int(tmp)
}

/**
*	@brief 目前是否处于连接中
 */
func (this *Player) IsConnect() bool {
	this.RLock()
	defer this.RUnlock()
	if this.m_pCon == nil {
		return false
	}
	return true
}

/**
*	@brief 断开连接
 */
func (this *Player) DisConnect() bool {
	this.Lock()
	defer this.Unlock()
	if this.m_pCon == nil {
		return false
	}
	err := this.m_pCon.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Idx=%d, %s", this.m_Idx, err.Error())
	}
	this.m_pCon = nil
	return true
}

func (this *Player) _open_channel() { //打开消息管道
	if this.m_Channel != nil {
		return
	}
	this.m_Channel_Closed = false
	this.m_Channel = make(chan interface{}, 10) //10个通道的最大数量
	go func() {
		for {
			//只要消息管道没有关闭，就一直等待消息
			v, open := <-this.m_Channel
			if !open {
				this.m_Channel = nil
				break
			}
			opType, _ := v.(int)
			this._accept_op(opType)
		}
		this._close_channel()
	}()
}

func (this *Player) _close_channel() { //关闭通道
	if this.m_Channel_Closed == true {
		return
	} //已经准备要关闭了
	if this.m_Channel == nil {
		return
	}
	close(this.m_Channel)
	this.m_Channel_Closed = true
}

func (this *Player) _accept_op(opType int) {
	switch opType {
	case base.OP_REGISTER:
		this._send_Message(message.BuildLoginRegister(this.m_Idx, this.Username, this.Password))
		break
	case base.OP_LOGIN:
		this._send_Message(message.BuildLoginLogin(this.m_Idx, this.Username, this.Password))
		break
	case base.OP_SYNC:
		//this._send_Message(BuildGetStarInfo(this.m_Idx))
		//this.syncData()
		//go func() {
		//	for {
		//		if this.Game {
		//			this.seq++
		//			sessionID := this.seq * 10000  + this.m_Idx
		//			this._send_Message(BuildSyncData(sessionID))
		//		}
		//		time.Sleep(time.Duration(this.synctime) * time.Second)
		//	}
		//}()
		break
	}
}

//func (this *Player) syncData() {
//	if this.Game {
//		this.seq++
//		sessionID := this.seq*10000 + this.m_Idx
//		this.syncStartTime = time.Now()
//		this._send_Message(BuildSyncData(sessionID))
//	}
//}

//func (this *Player)_deal_register(){//进行登录测试
//	if this.m_State == base.PLAYER_STATE_NONE{//无状态, 好吧 那就尝试请求登录吧
//
//
//		//fmt.Printf("=================== Account = %s LoginSession\n", account);
//		//request.Req_Client_Login = &protocol.Req_Client_Login_{ Account: account),Platform:"00") }
//
//		this.m_State =  base.PLAYER_STATE_LOGIN;
//	}else if this.m_State == base.PLAYER_STATE_LOGIN{//登录状态中， 好吧先不进行任何处理
//
//	}else if this.m_State == base.PLAYER_STATE_GAME{//登录成功状态， 好吧进行相应处理
//		this.DisConnect();
//		this.m_State =  base.PLAYER_STATE_NONE;
//	}
//}

func (this *Player) _send_Message(message *protocol.C2GS) { //发送消息
	this.RLock()
	defer this.RUnlock()
	if this.m_pCon == nil {
		return
	}

	buff, err := proto.Marshal(message)

	if err != nil { //序列化失败
		fmt.Printf(">> send message err = %v\n", err.Error())
		return
	}
	//fmt.Printf("%v => send: %v %v\n", this.Username, message, time.Now())
	if this.isCipher() {
		buff = xxtea.Encrypt(buff, []byte(this.secretkey))
		//fmt.Printf("key %v", this.secretkey)
	}

	m := make([]byte, len(buff)+2)
	binary.LittleEndian.PutUint16(m, uint16(len(buff))) //
	copy(m[2:], buff)

	this.m_pCon.Write(m) //把消息发送出去
}

func (this *Player) _receive_Message(message *protocol.GS2C) int32{ //接收消息
	//code := message.GetSequence()
	fmt.Printf("%v => receive: %v %v\n", this.Username, message, time.Now())
	if message.GetLoginRegisterRet() != nil {
		this.handlerRegisterRet(message.GetLoginRegisterRet())
		return 0
	}
	if message.GetLoginLoginRet() != nil {
		this.handlerLoginRet(message.GetLoginLoginRet())
		return 0
	}

	if message.GetResultPush() != nil {
		log.Debug("player %v test over, result : %v",  this.Uid, message.GetResultPush())
		return RET_TYPE_RESULT
	}

	if message.GetGetStarsSelectRet() != nil {
		this.handlerGetStarsSelectRet(message.GetGetStarsSelectRet())
		return 0
	}

	if message.GetSelectStarRet() != nil {
		this.handlerSelectStarRet(message.GetSelectStarRet())
	}
	if message.GetGetStarInfoRet() != nil {
		if this.handlerGetStarInfoRet(message.GetGetStarInfoRet()) {
			return 0
		}
	}

	if message.GetLoginServerRet() != nil {
		this.Game = true
		this.handlerLoginServerRet(message.GetLoginServerRet())
		return 0
	} else if message.GetKickoffPush() != nil {
		//this.handlerKickoffPush(message.GetKickoffPush())
	}

	this.AcceptResult(message)
	nextMessage := this.NextMessage()
	if nextMessage == nil {
		log.Debug("player %v test done!", this.Uid)
		return RET_TYPE_DONE
	}
	this._send_Message(nextMessage)
	return 0
}

func (this *Player) handlerRegisterRet(request *protocol.LoginRegisterRet) { //接收消息
	if request.GetResult() != protocol.Register_Result_registerSuccess {
		return
	}
	this.Uid = request.GetUid()
	this.Token = request.GetToken()
	this.ReconnectGame()
	this._send_Message(message.BuildLoginServer(this.m_Idx, this.Uid, this.Token))
}

//func (this *Player) handleMessage() { //接收消息
//	//go func() {
//	//	time.Sleep(time.Duration(this.synctime) * time.Second)
//	//	this._send_Message(BuildGetStarInfo(this.m_Idx))
//	//} ()
//	this._send_Message(BuildGetStarInfo(this.m_Idx))
//}

func (this *Player) handlerLoginRet(request *protocol.LoginLoginRet) { //接收消息
	if request.GetResult() != protocol.Login_Result_loginSuccess {
		return
	}
	this.Uid = request.GetUid()
	this.Token = request.GetToken()
	this.ReconnectGame()
	this._send_Message(message.BuildLoginServer(this.m_Idx, this.Uid, this.Token))
}

func (this *Player) handlerKickoffPush(request *protocol.KickoffPush) { //接收消息
	//1用户在其他地方登陆 2长时间没操作 3服务器关闭
	if request.GetType() == 3 {
		return
	}
	//被踢下线,等待5秒自动重连
	time.Sleep(5 * time.Second)
	this.ReconnectPassport()
	this._send_Message(message.BuildLoginLogin(this.m_Idx, this.Username, this.Password))
}

//func (this *Player) handlerCreateRoleRet(request *protocol.CreateRoleRet) {
//	if request.GetResult() != 0 {
//		return
//	}
//	this._send_Message(BuildCreateRole(this.m_Idx,))
//}

//func (this *Player) handlerJoinGameRet(request *protocol.JoinGameRet) {
//	if request
//}

func (this *Player) handlerLeaveGameRet(ret *protocol.LeaveGameRet) {
	if ret.Result != true {
		return
	}
	this._send_Message(message.BuildLeaveGame(this.m_Idx,this.Uid))
}

func (this *Player) handlerLoginServerRet(ret *protocol.LoginSeverRet) {
	this._send_Message(message.BuildGetStarInfo())
}

func (this *Player) handlerGetStarInfoRet(ret *protocol.GetStarInfoRet) bool{
	if ret.CurrentStar == nil {
		this._send_Message(message.BuildGetStarsSelect(-1))
		return true
	} else {
		this.TestSuite()
		return false
	}
}

func (this *Player) handlerGetStarsSelectRet(ret *protocol.GetStarsSelectRet) {
	var selectType int32 = 1
	this.starType = selectType
	this._send_Message(message.BuildSelectStar(selectType))
}

func (this *Player) handlerSelectStarRet(ret *protocol.SelectStarRet) {
	this.starID = ret.CurrentStar.StarID
	this.TestSuite()
}
