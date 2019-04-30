/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/11/15
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package http

import (
	"net/http"
	"gok/module/center/conf"
	"aliens/log"
	"github.com/name5566/leaf/chanrpc"
	"aliens/common/helper"
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
	RESULT_CODE_USER_EXIST	int = 6 //用户已存在
	RESULT_CODE_FAILED		  int = 7 //失败
)

const (
	PARAM_SIGN = "sign"
)

func Init(chanRpc *chanrpc.Server) {
	httpServer := http.NewServeMux()

	httpProxy := &helper.HttpProxy{}
	httpProxy.InitMux(chanRpc, httpServer)

	//httpProxy.RegisterFunc("/api/user/loginoauth", ChannelLogin)
	httpProxy.RegisterFunc("/api/user/pay", Pay)

	httpServer.HandleFunc("/api/user/mail", Mail)
	//httpProxy.RegisterFunc("/api/center", Center) //更新维护状态

	//http.HandleFunc("/api/user/mail", Mail)

	//关闭GC
	//debug.SetGCPercent(-1)
	//运行trace
	//http.HandleFunc("/start", traces)
	//停止trace
	//http.HandleFunc("/stop", traceStop)
	//手动GC

	//中控接口和渠道登录接口只是转发，不需要串行化tcp_conn
	httpServer.HandleFunc("/api/center", Center) //更新维护状态
	go func() {
		log.Debug("CenterHTTP:%v", conf.Server.HTTPAddress)
		log.Fatal("%v", http.ListenAndServe(conf.Server.HTTPAddress, httpServer))
	}()

}

func Close() {

}

//手动GC
//func gc(w http.ResponseWriter, r *http.Request) {
//	runtime.GC()
//	w.Write([]byte("StartGC"))
//}

//运行trace
//func traces(w http.ResponseWriter, r *http.Request){
//	f, err := os.Create("trace.out")
//	if err != nil {
//		panic(err)
//	}
//
//
//	err = trace.Start(f)
//	if err != nil {
//		panic(err)
//	}
//	w.Write([]byte("TrancStart"))
//	fmt.Println("StartTrancs")
//}

//停止trace
//func traceStop(w http.ResponseWriter, r *http.Request){
//	trace.Stop()
//	w.Write([]byte("TrancStop"))
//	fmt.Println("StopTrancs")
//}

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
