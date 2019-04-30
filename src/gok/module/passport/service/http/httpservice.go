package http

import (
	"github.com/name5566/leaf/chanrpc"
	"aliens/log"
	"net/http"
	"gok/module/passport/conf"
	"github.com/gogo/protobuf/sortkeys"
	"aliens/common/cipher"
)

const (
	RESULT_CODE_SUCCESS       int = 0  //成功
	RESULT_CODE_INVALID_PARAM int = 1  //无效的参数
	RESULT_CODE_INVALID_SIGN  int = 2  //签名失败
	RESULT_CODE_USER_NOTFOUND  int = 3 //用户未找到
	RESULT_CODE_DATABASE_EXCEPTION  int = 4 //数据库操作异常
	RESULT_CODE_SERVER_EXCEPTION  int = 5 //服务器内部处理异常
)

const (
	PARAM_SIGN = "sign"
)

func Init(chanRpc *chanrpc.Server) {
	httpServer := http.NewServeMux()

	//httpProxy := &helper.HttpProxy{}
	//httpProxy.InitMux(chanRpc, httpServer)

	httpServer.HandleFunc("/api/notice", Notice) //公告信息

	httpServer.HandleFunc("/api/wx", Wx) //

	httpServer.HandleFunc("/api/flag", Flag) //

	httpServer.HandleFunc("/api/wechat", Wechat)

	httpServer.HandleFunc("/api/user/loginoauth", ChannelLogin)

	httpServer.HandleFunc("/api/vivo/vivoLogin", VivoLogin)

	httpServer.HandleFunc("/api/vivo/vivoToken", VivoGetToken)

	httpServer.HandleFunc("/api/vivo/vivoRefresh", VivoRefresh)

	httpServer.HandleFunc("/api/test/custompush", CustomPush)

	httpServer.HandleFunc("/api/wx/refresh", Refresh)

	//httpServer.HandleFunc("/api/wx/daygift", DayGift)

	go func() {
		log.Debug("PassportHTTP:%v", conf.Server.HTTPAddress)
		log.Fatal("%v", http.ListenAndServe(conf.Server.HTTPAddress, httpServer))
	}()
}

func Close() {

}

func isSignSuccess(reqData map[string][]string) bool {
	if !conf.Server.IsSign {
		return true
	}

	var signText string
	var signData string
	var strKeys = []string{}
	for key := range reqData {
		strKeys = append(strKeys, key)
	}
	sortkeys.Strings(strKeys)

	for _, value := range strKeys {
		if value != "sign" {
			for _,value := range reqData[value] {
				signText += value
			}
		}
	}
	signText += conf.Server.AppKey
	signResult := cipher.MD5Hash(signText)
	for _,value := range reqData[PARAM_SIGN]{
		signData += value
	}
	if signResult != signData {
		log.Debug("signResult %v : sign : %v", signResult, signData)
		return false
	}

	return true
}