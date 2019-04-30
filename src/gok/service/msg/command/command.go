/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/4/11
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package command

//系统消息
type Command int32

type SystemMessage struct {
	Cmd 	Command		//系统消息编号
	Param  []interface{}	//系统消息参数
}

func NewSystemMessage(cmd Command, param ...interface{}) *SystemMessage {
	return &SystemMessage{
		Cmd:cmd,
		Param:param,
	}
}

const (
	//LOGIN     //用户登录初始化
	//LOGOUT 		    	//用户登出
	//RELEASE	Command = iota		//释放用户内存
	TIMER   Command = iota      //处理用户定时处理
	//LOGINPUSH //登录推送
	MESSAGEPUSH //消息推送
	//EVENT
	//KICK         //T人和释放内存

	AGENT_CLOSE  //网络连接断开

)

