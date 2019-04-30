/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/11/10
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package message

import (
	"gok/service/msg/protocol"
	"aliens/common/character"
	"time"
)

//构建注册消息
//登录服务器注册账号
//message login_register{
//optional string username = 1;	    //用户名
//optional string password = 2;		//密码
//}
//
////登录服务器注册账号返回
//message login_register_ret {
//optional register_Result result = 1;    //注册结果
//optional int64 uid = 2;                 //用户id 注册成功返回此字段
//optional string token = 3;              //登录令牌 注册成功返回此字段
//optional string gameServer = 4;         //返回游戏服务器地址
//optional string msg = 5;                //反馈消息 登录失败返回
//}

//func BuildSyncData(session int32) *protocol.C2GS {
//	message := &protocol.C2GS{
//		Session:  session),
//		Sequence: []int32{20},
//	}
//	resources := []*protocol.DataResource{}
//	resources = append(resources, &protocol.DataResource{
//		Id:    1000),
//		Value: 100),
//		Limit: 200),
//	})
//	request := &protocol.SyncData{
//		Resource: resources,
//	}
//
//	message.SyncData = request
//	return message
//}

func BuildLoginRegister(session int32, username string, passwd string) *protocol.C2GS {
	message := &protocol.C2GS{
		Session:  session,
		Sequence: []int32{6},
	}

	request := &protocol.LoginRegister{
		Username: username,
		Password: passwd,
	}

	message.LoginRegister = request
	return message
}

func BuildLoginLogin(session int32, username string, passwd string) *protocol.C2GS {
	message := &protocol.C2GS{
		Session:  session,
		Sequence: []int32{7},
	}

	request := &protocol.LoginLogin{
		Username: username,
		Password: passwd,
	}

	message.LoginLogin = request
	return message
}

func BuildLoginServer(session int32, uid int32, token string) *protocol.C2GS {
	message := &protocol.C2GS{
		Session:  session,
		Sequence: []int32{11},
	}

	request := &protocol.LoginServer{
		Token:   token,
		UserId:  uid,
		Version: "8.88.122201",
		InviteType:&protocol.InviteType{Type:0, RefType:0, BuildType:0},
	}

	message.LoginServer = request
	return message
}

func BuildGetStarInfo() *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: []int32{20},
	}
	message.GetStarInfo = &protocol.GetStarInfo{}
	return message
}

func int32AddZeroToString(num int32) string{
	str := character.Int32ToString(num)
	if num < 10 {
		str = "0" + str
	}
	return str
}




func BuildGetStarShield() *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: []int32{303},
	}
	message.GetStarShield = &protocol.GetStarShield{}
	return message
}

func BuildGetStarHistory(starID int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: []int32{380},
	}
	request := &protocol.GetStarHistory{
		StarID:  starID,
	}
	message.GetStarHistory = request
	return message
}

func BuildGetStarStatistics(uid int32, starID int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: []int32{371},
	}
	request := &protocol.GetStarStatistics{
		Uid:	uid,
		StarID:	starID,
	}
	message.GetStarStatistics = request
	return message
}


func BuildGetNewsFeedList(uid int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{261},
	}
	request := &protocol.GetNewsfeedList{
		Id:			uid,
	}
	message.GetNewsFeedList = request
	return message
}

func BuildSaleInfo(uid int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{253},
	}
	request := &protocol.GetSaleInfo{
		Id:	uid,
	}
	message.SaleInfo = request
	return message
}

func BuildAutoAddBeliever(times int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{402},
	}
	request := &protocol.AutoAddBeliever{
		Times: 0,
	}
	message.AutoAddBeliever = request
	return message
}

func buildMessage(session int32, portID int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Session: session,
		Sequence: []int32{portID},
	}
	return message
}

// 用户模块

//func BuildCreateRole(session int32, icon int32, nickname string) *protocol.C2GS {
//	message := buildMessage(session, 12)
//	message.CreateRole = &protocol.CreateRole{
//		Icon: icon,
//		Nickname: nickname,
//	}
//	return message
//}

//func BuildJoinGame(session int32, uid int32) * protocol.C2GS {
//	message := buildMessage(session, 13)
//	message.JoinGame = &protocol.JoinGame{
//		Id: uid,
//	}
//	return message
//}

func BuildLeaveGame(session int32, uid int32) *protocol.C2GS {
	message := buildMessage(session, 14)
	message.LeaveGame = &protocol.LeaveGame{}
	return message
}

//func BuildBuildStarBuilding(session int32, starID int32, buildingType int32, level int32,) *protocol.C2GS {
//	message := buildMessage(session, 21)
//	message.BuildStarBuilding = &protocol.BuildStarBuilding{
//		StarID: starID,
//		BuildingType: buildingType,
//		Level: level,
//		BelieverId: []string{},
//		//Faith: ,
//		//Guide: ,
//	}
//	return message
//}

//-----------------------------------------internal_add_attach--------------------------------------
func BuildAddAttach( starType int32, faith int32, power int32, diamond int32, gayPoint int32, isItem bool, ID int32, itemNum int32, isBeliever bool, level int32, believerNum int32) *protocol.C2GS{

	items := []*protocol.BagItem{}
	believers := []*protocol.BelieverInfo{}
	if isItem {
		//itemID := 10000 + starType*100
		//var itemNum int32 = 1
		for i := 1; i <= int(ID); i++ {
			//itemID += int32AddZeroToString(int32(i))
			itemID := 10000 + starType*100 + int32(i)
			item := &protocol.BagItem{Id:itemID,Num:itemNum}
			items = append(items, item)
		}
	}
	if isBeliever {
		typeStr := int32AddZeroToString(starType)
		believerID := "b" + typeStr + character.Int32ToString(level) + "1"
		believer := &protocol.BelieverInfo{
			Id:believerID,
			Num:believerNum,
		}
		believers = append(believers, believer)
	}

	message := &protocol.C2GS{
		Sequence:	[]int32{666},
	}
	request := &protocol.AddAttach{
		Faith: faith,
		Power: power,
		Diamond: diamond,
		GayPoint: gayPoint,
		Items: items,
		Believers: believers,
	}
	message.AddAttach = request
	return message
}

func BuildAddFaith(faith int32) *protocol.C2GS{
	message := &protocol.C2GS{
		Sequence:	[]int32{666},
	}
	request := &protocol.AddAttach{
		Faith: faith,
	}
	message.AddAttach = request
	return message
}
func BuildAddPower(power int32) *protocol.C2GS{
	message := &protocol.C2GS{
		Sequence:	[]int32{666},
	}
	request := &protocol.AddAttach{
		Power: power,
	}
	message.AddAttach = request
	return message
}
func BuildAddDiamond(diamond int32) *protocol.C2GS{
	message := &protocol.C2GS{
		Sequence:	[]int32{666},
	}
	request := &protocol.AddAttach{
		Diamond:diamond,
	}
	message.AddAttach = request
	return message
}
func BuildAddGayPoint(gayPoint int32) *protocol.C2GS{
	message := &protocol.C2GS{
		Sequence:	[]int32{666},
	}
	request := &protocol.AddAttach{
		GayPoint:gayPoint,
	}
	message.AddAttach = request
	return message
}
func BuildAddItems( starType int32, ID int32, itemNum int32 ) *protocol.C2GS{
	message := &protocol.C2GS{
		Sequence:	[]int32{666},
	}
	items := []*protocol.BagItem{}
	for i := 1; i <= int(ID); i++ {
		itemID := 10000 + starType*100 + ID
		item := &protocol.BagItem{Id:itemID,Num:itemNum}
		items = append(items, item)
	}
	request := &protocol.AddAttach{
		Items:items,
	}
	message.AddAttach = request
	return message
}
func BuildAddBelievers( starType int32, level int32, believerNum int32) *protocol.C2GS{
	message := &protocol.C2GS{
		Sequence:	[]int32{666},
	}
	believers := []*protocol.BelieverInfo{}
	typeStr := int32AddZeroToString(starType)
	believerID := "b" + typeStr + character.Int32ToString(level) + "1"
	believer := &protocol.BelieverInfo{
		Id:believerID,
		Num:believerNum,
	}
	believers = append(believers, believer)
	request := &protocol.AddAttach{
		Believers:believers,
	}
	message.AddAttach = request
	return message
}

//--------------------------------------star module----------------------------------------------
func BuildUpgradeBeliever(uid int32, starType int32) *protocol.C2GS {
	typeStr := character.Int32ToString(starType)
	if starType < 10 {
		typeStr = "0" + typeStr
	}
	message := &protocol.C2GS{
		Sequence:	[]int32{400},
	}
	request := &protocol.UpgradeBeliever{
		SelectID: 	"b"+ typeStr + "11",
		MatchID: 	"b"+ typeStr + "11",
		Uid:		uid,
		Faith:		0,
	}
	message.UpgradeBeliever = request
	return message
}

func BuildBuildStarBuilding(starID int32, buildingType int32, level int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{21},
	}
	request := &protocol.BuildStarBuilding{
		StarID: starID,
		BuildingType: buildingType,
		Level: level,
		BelieverId: []string{},
		//Faith: ,
		//Guide: ,
	}
	message.BuildStarBuilding = request
	return message
}

func BuildUpdateStarBuildEnd(uid int32, buildingType int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{350	},
	}
	request := &protocol.UpdateStarBuildEnd{
		Uid:			uid,
		BuildingType:	buildingType,
	}
	message.UpdateStarBuildEnd = request
	return message
}

func BuildAccUpdateStarBuild(uid int32, buildingType int32, starType int32, level int32) *protocol.C2GS{
	typeStr := int32AddZeroToString(starType)

	believerID := []string{}
	believerNum := 0
	believer := "b" + typeStr + "61"
	if level == 1 {
		believerNum = 1
	} else if level == 2 {
		believerNum = 1
	} else if level == 3 {
		believerNum = 3
	} else if level == 4 {
		believerNum = 9
	} else if level == 5 {
		believerNum = 20
	}
	for i := 0; i < believerNum; i++ {
		believerID = append(believerID, believer)
	}
	message := &protocol.C2GS{
		Sequence: 	[]int32{351},
	}
	request := &protocol.AccUpdateStarBuild{
		Uid:			uid,
		BuildingType: 	buildingType,
		BelieverId: 	believerID,
	}
	message.AccUpdateStarBuild = request
	return message
}

func BuildGetBuildingFaith(uid int32, buildingType int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{355},
	}
	request := &protocol.GetBuildingFaith{
		Uid:			uid,
		BuildingType:	buildingType,
	}
	message.GetBuildingFaith = request
	return message
}

func BuildGetItemGroup(uid int32, starType int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{451},
	}
	request := &protocol.GetItemGroup{
		Uid:		uid,
		StarType:	starType,
	}
	message.GetItemGroup = request
	return message
}

func BuildActiveGroup(uid int32, starType int32, groupID int32) *protocol.C2GS {
	activeGroupID := 10000 + starType*100 + int32(groupID)
	if starType == 2 && groupID > 5 {
		activeGroupID += 1
	}
	itemIDs := []int32{}
	for i := 2; i <= 6 ; i++ {
		itemID := 10000+starType*100+int32(i)
		itemIDs = append(itemIDs, itemID)
		if i ==4 && groupID <= 5{
			break
		}
	}
	message := &protocol.C2GS{
		Sequence: []int32{453},
	}
	request := &protocol.ActiveGroup{
		Uid:uid,
		GroupID:activeGroupID,
		ItemID:itemIDs,
	}
	message.ActiveGroup = request
	return message
}

func BuildDrawCivilizationReward(starID int32, level int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: []int32{369},
	}
	request := &protocol.DrawCivilizationReward{
		StarID:starID,
		DrawLevel:level,
	}
	message.DrawCivilizationReward = request
	return message
}

func BuildSetBuildings(starID int32, totalLevel int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: []int32{667},
	}
	request := &protocol.SetBuildings{
		StarID: starID,
		Level: totalLevel,
	}
	message.SetBuildings = request
	return  message
}

func BuildSetBelievers(totalLevel int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: []int32{668},
	}
	request := &protocol.SetBelievers{
		Level:totalLevel,
	}
	message.SetBelievers = request
	return  message
}

//---------------------------------mail------------------------------------
func BuildCreateMail(uid int32, starType int32, faith string, power string, diamond string, gay_point string, itemNum string, believerLevel string, believerNum string ) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: []int32{73},
	}
	typeStr := int32AddZeroToString(starType)
	believerID := "b" + typeStr + believerLevel + "1"
	itemID := "1" + typeStr

	itemStr := ""
	for i := 1; i <= 5; i++ {
		itemStr += "{\"id\":" + itemID + int32AddZeroToString(int32(i)) + ",\"num\":"+itemNum+"}"
		if i != 5{
			itemStr += ","
		}
	}

	attach := "{\"diamond\":"+diamond+",\"gay_point\":"+gay_point+",\"faith\":"+faith+",\"power\":"+power+",\"item\":["+itemStr+"],\"believer\":[{\"id\":\""+ believerID +"\",\"num\":"+ believerNum +"}]}"
	request := &protocol.CreateMail{
		Uid: uid,
		Title: "",
		Content: "",
		MailAttach: attach,
	}
	message.CreateMail = request
	return message
}

func BuildGetAllMail(uid int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{70},
	}
	request := &protocol.GetAllMail{
		Uid:		uid,
	}
	message.GetAllMail = request
	return message
}

func BuildDrawMail(uid int32, maiID int64) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{71},
	}
	request := &protocol.DrawMail{
		Uid:uid,
		MailID:maiID,
	}
	message.DrawMail = request
	return message
}

func BuildRemoveMail(uid int32,mailID []int64) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{72},
	}
	request := &protocol.RemoveMail{
		Uid:uid,
		MailID:mailID,
	}
	message.RemoveMail = request
	return message
}

//-----------------------------------community-------------------------------
func BuildFollow(uid int32, followerID int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{220},
	}
	request := &protocol.Follow{
		Id:uid,
		FollowerID:followerID,
	}
	message.Follow = request
	return message
}

func BuildUnFollow(uid int32, unFollowerID int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{221},
	}
	request := &protocol.Unfollow{
		Id:uid,
		UnfollowerID:unFollowerID,
	}
	message.Unfollow = request
	return message
}
func BuildGetFollowerList(uid int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{222},
	}
	request := &protocol.GetFollowerList{
		Id:uid,
	}
	message.GetFollowerList = request
	return message
}
func BuildGetFollowingList(uid int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{224},
	}
	request := &protocol.GetFollowingList{
		Id:uid,
	}
	message.GetFollowingList = request
	return message
}

func BuildGetReceiveMoments(uid int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{240},
	}
	request := &protocol.GetReceiveMoments{
		Uid:uid,
		BeforeTime:time.Now().Unix(),
		Count:30,
		Offset:0,
	}
	message.GetReceiveMoments = request
	return message
}

func BuildGetPublicMoments(uid int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{241},
	}
	request := &protocol.GetPublicMoments{
		Uid:uid,
		BeforeTime:time.Now().Unix(),
		Count:30,
		Offset:0,
	}
	message.GetPublicMoments = request
	return message
}

func BuildPublicMoment(uid int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{243},
	}
	request := &protocol.PublicMoment{
		Uid:uid,
		Type:1,
		RefID:666,
	}
	message.PublicMoment = request
	return message
}

func BuildRemoveMoment(uid int32, momentID []string) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{244},
	}
	request := &protocol.RemoveMoments{
		Uid:uid,
		MomentsID:momentID,
		SaleID:uid,
	}
	message.RemoveMoments = request
	return message
}

//------------------------------------------trade----------------------------------------
func concatItemID(starType int32, itemID int32) int32 {
	return 	10000 + starType*100 + itemID
}

func BuildGetGoodsInfo(uid int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{470},
	}
	request := &protocol.GetGoodsInfo{
		Uid:uid,
	}
	message.GetGoodsInfo = request
	return message
}

//TODO
//func BuildBuyGoods() *protocol.C2GS {
//	message := &protocol.C2GS{
//		Sequence:	[]int32{471},
//	}
//	request := &protocol.BuyGoods{
//
//	}
//	message.BuyGoods = request
//	return message
//}

func BuildPublicGoods(uid int32, starType int32, ID int32) *protocol.C2GS {
	itemID := concatItemID(starType, ID)
	message := &protocol.C2GS{
		Sequence:	[]int32{472},
	}
	request := &protocol.PublicGoods{
		Uid:uid,
		Goods:&protocol.Goods{
			Id:itemID,
			Num:1,
			Price:100,
		},
	}
	message.PublicGoods = request
	return message
}

func BuildCancelGoods(uid int32, starType int32, ID int32) *protocol.C2GS {
	itemID := concatItemID(starType, ID)
	message := &protocol.C2GS{
		Sequence:	[]int32{473},
	}
	request := &protocol.CancelGoods{
		Uid:uid,
		Itemid:itemID,
	}
	message.CancelGoods = request
	return message
}

func BuildAddSale(uid int32, starType int32, ID int32) *protocol.C2GS {
	itemID := concatItemID(starType, ID)
	message := &protocol.C2GS{
		Sequence:	[]int32{530},
	}
	request := &protocol.AddSale{
		Id:uid,
		ItemID:itemID,
	}
	message.AddSale = request
	return message
}

func BuildRemoveSale(uid int32, starType int32, ID int32) *protocol.C2GS {
	itemID := concatItemID(starType, ID)
	message := &protocol.C2GS{
		Sequence:	[]int32{531},
	}
	request := &protocol.RemoveSale{
		Id:uid,
		ItemID:itemID,
	}
	message.RemoveSale = request
	return message
}

func BuildGetSale(uid int32) *protocol.C2GS {//获取指定uid的sale
	message := &protocol.C2GS{
		Sequence:	[]int32{532},
	}
	request := &protocol.GetSale{
		Id:uid,
	}
	message.GetSale = request
	return message
}

func BuildGetSales(uid []int32) *protocol.C2GS {//获取一组uid的sale
	message := &protocol.C2GS{
		Sequence:	[]int32{533},
	}
	request := &protocol.GetSales{
		Id:uid,
	}
	message.GetSales = request
	return message
}

//----------------------------------dial-------------------------------------
func BuildGetFaith(starType int32) *protocol.C2GS {
	typeStr := character.Int32ToString(starType)
	if starType < 10 {
		typeStr = "0" + typeStr
	}
	believerID := "b" + typeStr + "11"
	message := &protocol.C2GS{
		Sequence: 	[]int32{65},
	}
	request := &protocol.GetFaith{
		EventID:0,
		BelieverID:[]string{believerID,believerID,believerID},
	}
	message.GetFaith = request
	return message
}

func BuildGetBeliever(eventID int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: 	[]int32{66},
	}
	request := &protocol.GetBeliever{
		EventID:eventID,
	}
	message.GetBeliever = request
	return message
}

func BuildLootFaith(eventID int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: 	[]int32{67},
	}
	request := &protocol.LootFaith{
		EventID:eventID,
		Faith:200,
	}
	message.LootFaith = request
	return message
}

func BuildAtkBuilding(eventID int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: 	[]int32{68},
	}
	request := &protocol.AtkStarBuilding{
		EventID:eventID,
		BuildingID:0,
		Success:false,
	}
	message.AtkStarBuilding = request
	return message
}

func BuildLootBeliever(eventID int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: 	[]int32{69},
	}
	request := &protocol.LootBeliever{
		EventID:eventID,
	}
	message.LootBeliever = request
	return message
}


func BuildRandomDial() *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{106},
	}
	request := &protocol.RandomDial{}
	message.RandomDial = request
	return message
}

func BuildRandomTarget(uid int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: 	[]int32{40},
	}
	request := &protocol.RandomTarget{
		Uid:uid,
		EventID:uid,
	}
	message.RandomTarget = request
	return message
}

func BuildSelectEventTarget(uid int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: 	[]int32{34},
	}
	request := &protocol.SelectEventTarget{
		EventID:uid,
		TargetId:0,
		Nickname:"",
	}
	message.SelectEventTarget = request
	return message
}

func BuildIntoEvent(uid int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: 	[]int32{30},
	}
	request := &protocol.IntoEvent{
		Uid:uid,
		EventID:uid,
	}
	message.IntoEvent = request
	return message
}

func BuildDoneEventStep(eventID int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{46},
	}
	request := &protocol.DoneEventStep{
		EventID:eventID,
		Step:2,
	}
	message.DoneEventStep = request
	return message
}
//-----------------------------------user------------------------------------
func BuildGetOnNotices(status int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:	[]int32{700},
	}
	request := &protocol.GetOnNotices{
		Status:status,
	}
	message.GetOnNotices = request
	return message
}

func BuildGetAvatar(uid int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: 	[]int32{15},
	}
	request := &protocol.GetAvatar{
		Uid:[]int32{uid},
	}
	message.GetAvatar = request
	return message
}

func BuildChangeDesc(desc string) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: 	[]int32{16},
	}
	message.ChangeDesc = desc
	return message
}

func BuildDisplayInfo(min int32, max int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: 	[]int32{22},
	}
	request := &protocol.RoleDisplayInfo{
		Min:min,
		Max:max,
	}
	message.RoleDisplayInfo = request
	return message
}

func BuildUpdateDisplay(uid int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: 	[]int32{23},
	}
	request := &protocol.UpdateDisplay{
		Id:uid,
	}
	message.UpdateDisplay = request
	return message
}

func BuildRankInfo(uid int32, rankType int32)  *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: 	[]int32{150},
	}
	request := &protocol.GetRankInfo{
		Uid:uid,
		Type:rankType,
	}
	message.GetRankInfo = request
	return message
}

func BuildRoleFalgInfo()  *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: 	[]int32{26},
	}
	request := &protocol.RoleFlagInfo{
	}
	message.RoleFlagInfo = request
	return message
}

func BuildUpdateRoleFlag(flag int32, value int32)  *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: 	[]int32{27},
	}
	request := &protocol.UpdateFlag{
		Flag:flag,
		Value:value,
	}
	message.UpdateFlag = request
	return message
}

func BuildUpdatePower() *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: 		[]int32{28},
		UpdatePower: 	&protocol.UpdatePower{},
	}
	return message
}

func BuildBelieverFlagInfo() *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:			[]int32{401},
		BelieverFlagInfo:	&protocol.BelieverFlagInfo{},
	}
	return message
}

func BuildGetItem(uid int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:			[]int32{450},
		GetBagItem:			&protocol.GetBagItem{Uid:uid},
	}
	return message
}

func BuildSearchUser(uid int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:			[]int32{236},
		SearchUser:			&protocol.SearchUser{Id:uid},
	}
	return message
}

func BuildDealList(uid int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:			[]int32{266},
		GetDealList:		&protocol.GetDealList{Id:uid},
	}
	return message
}

func BuildGlobalMessage() *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:			[]int32{281},
		GetGlobalMessage:	true,
	}
	return message
}

func BuildPublicShare(publicType int32, refID string) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:		[]int32{290},
		PublicShare:	&protocol.PublicShare{Type:publicType,RefID:refID},
	}
	return message
}

func BuildPublicWechatShare() *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:			[]int32{291},
		PublicWechatShare:	&protocol.PublicWechatShare{},
	}
	return message
}

func BuildRequestItem(searchID string) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:		[]int32{262},
		RequestItem:	&protocol.RequestItem{SearchID:searchID},
	}
	return message
}

func BuildAcceptItem(dealID string) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:			[]int32{264},
		AcceptItemRequest:	&protocol.AcceptItemRequest{DealID:dealID},
	}
	return message
}

func BuildRejectItem(dealID string) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:			[]int32{265},
		RejectItemRequest:	&protocol.RejectItemRequest{DealID:dealID},
	}
	return message
}

func BuildGetStarsSelect(num int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence:			[]int32{18},
		GetStarsSelect: 	&protocol.GetStarsSelect{Num:num},
	}
	return message
}

func BuildSelectStar(starType int32) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: 			[]int32{19},
		SelectStar:  		&protocol.SelectStar{StarType:starType},
	}
	return message
}
