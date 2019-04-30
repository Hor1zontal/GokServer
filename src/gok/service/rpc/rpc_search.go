/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/4/17
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package rpc

import (
	"gok/constant"
	"gok/service"
	"gok/service/msg/protocol"
)

var SearchServiceProxy = &searchHandler{&rpcHandler{serviceType:service.SERVICE_SEARCH_RPC}}

var uploadData = &protocol.C2GS{Sequence: []int32{540}, UpdateSearchData: &protocol.UpdateSearchData{Opts:make(map[int32]*protocol.StarOpt)}}

const (
	maxUploadLimit = 200 //星球索引数据量超过200 直接上传
)

type searchHandler struct {
	*rpcHandler
}

func (this *searchHandler) RandomTarget(target *protocol.RandomTarget) *protocol.RandomTargetRet {
	message := &protocol.C2GS{
		Sequence:     []int32{40},
		RandomTarget: target,
	}
	return this.HandleMessage(message).GetRandomTargetRet()
}

//同步更新搜索索引数据
func (this *searchHandler) RandomHelpTargets(uid int32, starType int32, count int32) *protocol.RandomHelpTargetRet {
	message := &protocol.C2GS{
		Sequence: []int32{542},
		RandomHelpTarget: &protocol.RandomHelpTarget{Uid:uid, StarType:starType, Count:count},
	}
	return this.HandleMessage(message).GetRandomHelpTargetRet()
}

func (this *searchHandler) UpdateRandomStar(starID int32) {
	opt := make(map[int32]int32)
	opt[constant.SEARCH_OPT_UPDATE_BELIEVER] = 0
	opt[constant.SEARCH_OPT_UPDATE_BUILDING] = 0
	message := &protocol.C2GS{
		Sequence: []int32{543},
		UpdateRandomStar: &protocol.UpdateRandomStar{StarID:starID, Opt:opt},
	}
	this.AsyncHandleMessage(message)
}


func (this *searchHandler) UpdateData(opt int32, starID int32, updateData int32) {
	opts := uploadData.UpdateSearchData.Opts[starID]
	if opts == nil {
		opts = &protocol.StarOpt{Opt:make(map[int32]int32)}
		uploadData.UpdateSearchData.Opts[starID] = opts
	}
	opts.Opt[opt] = updateData
	if len(opts.Opt) > maxUploadLimit {
		this.SyncSearchData()
	}
}

//同步更新搜索索引数据
func (this *searchHandler) SyncSearchData() {
	//没有更新数据不需要上传同步
	if len(uploadData.UpdateSearchData.Opts) == 0 {
		return
	}
	this.AsyncHandleMessage(uploadData)
	uploadData.UpdateSearchData.Opts = make(map[int32]*protocol.StarOpt)
}


//同步更新搜索索引数据
func (this *searchHandler) UpdateHelpData(uid int32, opt int32, param int32) {
	message := this.BuildHelpDatas([]int32{uid}, opt, param, false)
	this.AsyncHandleMessage(message)
}

func (this *searchHandler) UpdateHelpDatas(uids []int32, opt int32, param int32) {
	message := this.BuildHelpDatas(uids, opt, param, false)
	this.AsyncHandleMessage(message)
}

func (this *searchHandler) BuildHelpDatas(uids []int32, opt int32, param int32, sync bool) *protocol.C2GS {
	message := &protocol.C2GS{
		Sequence: []int32{541},
		UpdateSearchHelpData: &protocol.UpdateSearchHelpData{Uid:uids, Opt:opt, Param:param, Sync: sync},
	}
	return message
}

//操作同步更新到其他服务器
func (this *searchHandler) SyncOtherSearchService(request *protocol.C2GS, ignoreNode string) {
	this.AsyncBroadcastAllRemoteIgnore(request, ignoreNode)
	//this.handlenodemessage.HandleMessage(request)
}