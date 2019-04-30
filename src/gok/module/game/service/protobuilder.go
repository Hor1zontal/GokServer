/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/4/27
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package service

import (
	"gok/service/msg/protocol"
	"gok/module/game/db"
	"gok/module/game/user"
	"gok/module/game/cache"
	clusterCache "gok/module/cluster/cache"
	"time"
)

func BuildTasks(this *user.Session) []*protocol.Task {
	result := []*protocol.Task{}
	tasks := this.GetTasks()
	if tasks == nil {
		return result
	}
	for _, dbTask := range this.GetTasks() {
		result = append(result, BuildTask(dbTask))
	}
	return result
}

func FillDetailInfo(userInfo *protocol.UserInfo) {
	//userInfo.Nickname = cache.UserCache.GetUserNickname(userInfo.GetId()))
	//userInfo.Avatar = cache.UserCache.GetUserAvatar(userInfo.GetId()))
	if clusterCache.Cluster.IsUserOnline(userInfo.GetId()) {
		userInfo.LastOnlineTime = 0
	} else {
		userInfo.LastOnlineTime = cache.UserCache.GetUserOnlineTimestamp(userInfo.GetId())
	}
	//userInfo.HelpItemID =
	userInfo.StarType = cache.StarCache.GetUserActiveStarType(userInfo.GetId())
	userInfo.Help = cache.TradeCache.ExistItemHelp(userInfo.GetId())
}

func BuildFlagInfoResponse(user *user.Session) *protocol.RoleFlagInfoRet {
	return &protocol.RoleFlagInfoRet{
		Flag:user.GetFlagProtocol(),
	}
}

func BuildBelieverFlagInfoResponse(user *user.Session) *protocol.BelieverFlagInfoRet {
	flag := []string{}
	value := []bool{}
	time := []int64{}
	for _, flagItem := range user.GetBelieverFlags() {
		flag = append(flag, flagItem.ID)
		value = append(value, flagItem.Value)
		time = append(time, flagItem.UpdateTime.Unix())
	}
	return &protocol.BelieverFlagInfoRet{
		Id:    flag,
		Value: value,
		Time: time,
	}
}

func BuildFlagInfoPush(flag *db.DBRoleFlag) *protocol.GS2C {
	return &protocol.GS2C{
		Sequence:[]int32{1002},
		RoleFlagPush: flag.BuildProtocol(),
	}
}

func AddFlagInfoPush(flag *db.DBRoleFlag, response *protocol.GS2C) {
	response.Sequence = append(response.Sequence, 1002)
	response.RoleFlagPush = flag.BuildProtocol()
}

func BuildTaskPush(task *db.DBRoleEventTask) *protocol.GS2C {
	return &protocol.GS2C{
		Sequence:[]int32{1007},
		TaskPush: BuildTask(task),
	}
}

func AddTaskPush(task *db.DBRoleEventTask, response *protocol.GS2C) {
	response.Sequence = append(response.Sequence, 1007)
	response.TaskPush = BuildTask(task)
}

func BuildTask(dbTask *db.DBRoleEventTask) *protocol.Task {
	return &protocol.Task{
		Id:       dbTask.ID,
		BaseID:   dbTask.Type,
		State:    int32(dbTask.State),
		EventID:  dbTask.RefID,
		EndingID: dbTask.EndingID,
		Reward: &protocol.Reward{
			Faith:    dbTask.RewardFaith,
			Believer: dbTask.RewardBeliever,
			ItemID: dbTask.RewardItem,
			GayPoint: dbTask.RewardGayPoint,
		},
	}
}

//func BuildAssistRequestPush(assist *db.DBRoleAssist) *protocol.GS2C {
//	return &protocol.GS2C{
//		Sequence:[]int32{1008},
//		AssistEventRequestPush: &protocol.AssistEventRequestPush{
//			Request: BuildAssistRequest(assist),
//		},
//	}
//}
//
//func BuildAssistRequest(assist *db.DBRoleAssist) *protocol.AssistEventRequest {
//	return &protocol.AssistEventRequest{
//		EventID:   assist.EventID,
//		Uid:       assist.Uid,
//		Nickname:  assist.NickName,
//		Msg:       assist.Msg,
//		Timestamp: assist.CreateTime.Unix(),
//	}
//}

//func BuildAssistRequests(this *user.Session) []*protocol.AssistEventRequest {
//	assists := this.GetAssists()
//	result := []*protocol.AssistEventRequest{}
//	for _, assist := range assists {
//		result = append(result, BuildAssistRequest(assist))
//	}
//	return result
//}

//构建角色信息
func BuildRoleInfo(this *user.Session) *protocol.RoleInfo {
	return &protocol.RoleInfo{
		Id:              this.GetID(),
//		Icon:            this.GetIcon()),
		Nickname:        this.GetNickName(),
		Desc: 			 this.GetDesc(),
		//Level:           this.GetLevel(),
		//Exp:             this.GetExp(),
		Power:           this.GetPower(),
		PowerLimit:      this.GetPowerLimit(),
		//Limit:           this.GetPowerLimit()),
		Faith:           this.GetFaith(),
		UpdateTimestamp: this.GetLastPowerTime().Unix(),
		Diamond:         this.GetDiamond(),
		Flag:            this.GetFlagProtocol(),
		GayPoint:        this.GetGayPoint(),
		Avatar:          cache.UserCache.GetUserAvatar(this.GetID()),
		Subscribe:       this.GetSubscribe(),
		//CivilizationLv: this.GetCivilizationLevel(),
		//CivilizationProgress: this.GetCivilizationValue(),
	}
}

//构建角色属性变更推送
func BuildRoleSocialPush(this *user.Session) *protocol.GS2C {
	return this.BuildRoleSocialPush()
}

func AddRoleSocialPush(this *user.Session, response *protocol.GS2C) {
	response.Sequence = append(response.Sequence, 1004)
	response.RoleInfoPush = &protocol.RoleInfoPush{
		//Level:    this.GetLevel(),
		//Exp:      this.GetExp(),
		Power:      this.GetPower(),
		PowerLimit: this.GetPowerLimit(),
		Faith:      this.GetFaith(),
		Diamond:    this.GetDiamond(),
		GayPoint:   this.GetGayPoint(),
	}
}


func BuildScorePush(statistic *db.DBStatistics) *protocol.GS2C {
	return &protocol.GS2C{
		Sequence:[]int32{1060},
		ScorePush: &protocol.ScorePush{
			Type:statistic.ID,
			Value:int64(statistic.Value),
		}}
}

func AddScorePush(statistic *db.DBStatistics, response *protocol.GS2C) {
	response.Sequence = append(response.Sequence, 1060)
	response.ScorePush =  &protocol.ScorePush{
		Type:statistic.ID,
		Value:int64(statistic.Value),
	}
}

func BuildStrangerPush(stranger *protocol.Stranger) *protocol.GS2C {
	return &protocol.GS2C{
		Sequence:[]int32{1051},
		StrangerPush: stranger,
	}
}


func BuildFollowPush(followID int32) *protocol.GS2C {
	return &protocol.GS2C{
		Sequence:[]int32{1053},
		FollowPush: followID,
	}
}

func BuildNewsFeedPush(newsFeed *protocol.NewsFeed) *protocol.GS2C {
	return user.BuildNewsFeedPush(newsFeed)
}

func AddNewsFeedPush(newsFeed *protocol.NewsFeed, response *protocol.GS2C) {

	response.Sequence = append(response.Sequence, 1052)
	response.NewsFeedPush = newsFeed
}


func BuildFriendInfo(id int32) *protocol.UserInfo {
	return &protocol.UserInfo{
		Id:       id,
		//Nickname: cache.UserCache.GetUserNickname(id)),
		//Avatar: cache.UserCache.GetUserAvatar(id)),
	}
}

func BuildStrangerMessage(uid int32, nickname string, strangerType int32) *protocol.C2GS {
	return &protocol.C2GS{
		Sequence:[]int32{522},
		AddStranger: BuildStranger(uid, nickname, strangerType),
	}
}

func BuildStranger(uid int32, nickname string, strangerType int32) *protocol.Stranger {
	return &protocol.Stranger{
		Id:uid,
		Type:strangerType,
		Time:time.Now().Unix(),
		Nickname:nickname,
		Avatar:cache.UserCache.GetUserAvatar(uid),
	}
}

func BuildNewsFeedMessage1(uid int32, newsFeedType int32, param1 int32, param2 int32, param3 int32, ext []string) *protocol.C2GS {
	return &protocol.C2GS{
		Sequence:[]int32{523},
		AddNewsFeed: BuildNewsFeed1(uid, newsFeedType, param1, param2, param3, ext),
	}
}

func BuildNewsFeedMessage(uid int32, newsFeedType int32, param1 int32, param2 int32, param3 int32) *protocol.C2GS {
	return &protocol.C2GS{
		Sequence:[]int32{523},
		AddNewsFeed: BuildNewsFeed(uid, newsFeedType, param1, param2, param3),
	}
}

func BuildNewsFeed1(uid int32, newsFeedType int32, param1 int32, param2 int32, param3 int32, ext []string) *protocol.NewsFeed {
	return user.BuildNewsFeed1(uid, newsFeedType, param1, param2, param3, ext)
}


func BuildNewsFeed(uid int32, newsFeedType int32, param1 int32, param2 int32, param3 int32) *protocol.NewsFeed {
	return user.BuildNewsFeed(uid, newsFeedType, param1, param2, param3)
}
