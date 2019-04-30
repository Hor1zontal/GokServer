package manager

import (
	"gok/module/game/db"
	"time"
	"gok/service/exception"
	"gok/module/game/conf"
	"gok/service/msg/protocol"
	"aliens/log"
)

//用户任务管理
type RoleEventTaskManager struct {
	uid int32
	task *db.DBRoleEventTask

	//tasks map[int32]*db.DBRoleEventTask //物品id和物品对象的映射表
	//mutualEvents []*db.DBMutualEvent //最近的交互事件
}

//初始化
func (this *RoleEventTaskManager) Init(role *db.DBRole) {
	this.uid = role.UserID
	//this.tasks = make(map[int32]*db.DBRoleEventTask)
	//for _, task := range role.Task {
	//	this.tasks[task.ID] = task
	//}
	//this.mutualEvents = role.MutualEvents
}

//更新内存
func (this *RoleEventTaskManager) Update(role *db.DBRole) {
	//role.Task = this.GetTasks()
	//role.MutualEvents = this.mutualEvents
}

func (this *RoleEventTaskManager) GetTasks() []*db.DBRoleEventTask {
	if this.task != nil {
		return []*db.DBRoleEventTask{this.task}
	} else {
		return nil
	}
	//for _, task := range this.tasks {
	//	result = append(result, task)
	//}
	//return result
}

//func (this *RoleEventTaskManager) AddMutualEvent(targetID int32, mType int32, updateTime time.Time) {
//	if time.Now().Sub(updateTime) > SEVEN_DAY {
//		return
//	}
//	for _, mutualEvent := range this.mutualEvents {
//		if mutualEvent.Uid == targetID {
//			mutualEvent.Type = mType
//			mutualEvent.Time = updateTime
//			return
//		}
//	}
//	if len(this.mutualEvents) == constant.MUTUALEVENT_LIMIT {
//		this.mutualEvents = this.mutualEvents[1:constant.MUTUALEVENT_LIMIT]
//	}
//	this.mutualEvents = append(this.mutualEvents, &db.DBMutualEvent{Uid: targetID, Type:mType, Time: updateTime})
//}


//更新规则随机规则
//func (this *RoleEventTaskManager) UpdateTargetRule(targetID int32) {
//	if targetID == 0 {
//		return
//	}
//	//从复仇列表中删除
//	for index, mutualEvent := range this.mutualEvents {
//		if mutualEvent.Uid == targetID {
//			nextIndex := index + 1
//			if len(this.mutualEvents) == nextIndex {
//				this.mutualEvents = this.mutualEvents[:index]
//			} else {
//				this.mutualEvents = append(this.mutualEvents[:index], this.mutualEvents[index+1:]...)
//			}
//		}
//	}
//}


//随机交互用户
//func (this *RoleEventTaskManager) RandomMutualTarget() int32 {
//	threeDayUser := []int32{}
//	sevenDayUser := []int32{}
//	newMutualEvents := []*db.DBMutualEvent{} //最近的交互事件
//	for _, mutualEvent := range this.mutualEvents {
//		duration := time.Now().Sub(mutualEvent.Time)
//		if duration <= THREE_DAY {
//			threeDayUser = append(threeDayUser, mutualEvent.Uid)
//			newMutualEvents = append(newMutualEvents, mutualEvent)
//		} else if duration <= SEVEN_DAY {
//			sevenDayUser = append(sevenDayUser, mutualEvent.Uid)
//			newMutualEvents = append(newMutualEvents, mutualEvent)
//		} else {
//			//清楚过期的数据
//			break
//		}
//	}
//	if len(threeDayUser) > 0 {
//		return threeDayUser[rand.Intn(len(threeDayUser))]
//	}
//	if len(sevenDayUser) > 0 {
//		return sevenDayUser[rand.Intn(len(sevenDayUser))]
//	}
//	this.mutualEvents = newMutualEvents
//	return 0
//}

//func (this *RoleEventTaskManager) RandomFriendTarget() int32 {
//	friends := cache.UserCache.GetFriendArray(this.uid)
//	friendLen := len(friends)
//	if friendLen == 0 {
//		return 0
//	}
//	friendRandomIndex := rand.Intn(friendLen)
//	currTimestamp := time.Now().Unix()
//	for i:=friendRandomIndex; i<friendLen; i ++ {
//		if this.matchFriend(this.uid, friends[i], currTimestamp) {
//			return friends[i]
//		}
//	}
//	for i:=0; i<friendRandomIndex; i ++ {
//		if this.matchFriend(this.uid, friends[i], currTimestamp) {
//			return friends[i]
//		}
//	}
//	return 0
//}

//好友是否在1小时内被搜索到过,以及在线时间在七天内
//func (this *RoleEventTaskManager) matchFriend(uid string, friendID int32, currTimestamp int64) bool {
//	if cache.UserCache.ExistMutual(uid, character.Int32ToString(friendID)) {
//		return false
//	}
//	return (currTimestamp - cache.UserCache.GetUserOnlineTimestamp(character.StringToInt32(uid))) < int64(7 * util.SECOEND_OF_DAY)
//}

//func (this *RoleEventTaskManager) FilterDealOverdue() int32 {
//	var total int32 = 0
//	for id, deal := range this.mutualEvents {
//		if (time.Now().Sub(deal.Time)) {
//			if (deal.GetType() == constant.NEWSFEED_TYPE_REQUEST_ITEM) {
//				total ++
//			}
//			delete(this.deals, id)
//		}
//	}
//	return total
//}


func (this *RoleEventTaskManager) HaveTask() bool {
	return this.task != nil
}

func (this *RoleEventTaskManager) CleanTask() {
	this.task = nil
	//for _, task := range this.tasks {
	//	//if task.RefID != 0 {
	//	//	user.RoleTempManager.DeleteEvent()
	//	//}
	//	delete(this.tasks, task.ID)
	//}
}

//获取任务的的任务类型
//func (this *RoleEventTaskManager) GetTaskTypes() *set.HashSet {
//	result := set.NewHashSet()
//	for _, task := range this.tasks {
//		result.Add(task.Type)
//	}
//	return result
//}

func (this *RoleEventTaskManager) NewTask(taskType int32, refID int32) *db.DBRoleEventTask {
	base := conf.GetTaskData(taskType)
	if base == nil {
		exception.GameException(exception.TASK_BASE_NOTFOUND)
	}

	this.task = &db.DBRoleEventTask{
		ID:this.genTaskID(),
		Type:taskType,
		RefID:refID,
		State:db.TASK_STATE_RUNNING,
		CreateTime:time.Now(),
	}
	//log.Debug("time: %v, userID:%v, newTaskID: %v, refID: %v", this.task.CreateTime, this.uid, this.task.ID, this.task.RefID)
	//this.tasks[task.ID] = task
	return this.task
}

//新建任务
//func (this *RoleEventTaskManager) NewTask(taskType int32, refID int32, initFields map[string]int32) *db.DBRoleEventTask {
//	base := conf.GetTaskData(taskType)
//	if (base == nil) {
//		exception.GameException(exception.TASK_BASE_NOTFOUND)
//	}
//	task := &db.DBRoleEventTask{
//		ID:this.genTaskID(),
//		Type:taskType,
//		RefID:refID,
//		State:db.TASK_STATE_RUNNING,
//		CreateTime:time.Now(),
//	}
//
//	//从配置表和外部参数加载任务当前信息
//	results := []*db.TaskField{}
//	for _, field := range base.Fields {
//		var initValue int32 = 0
//		if (initFields != nil) {
//			initValue = initFields[field.Field]
//		}
//		results = append(results, &db.TaskField{
//			Name:field.Field,
//			Threshold:field.Threshold,
//			Value:initValue,
//		})
//	}
//	task.Fields = results
//
//	this.tasks[task.ID] = task
//	return task
//}

func (this *RoleEventTaskManager) GetTask(taskID int32) *db.DBRoleEventTask {
	if this.task == nil || this.task.ID != taskID {
		return nil
	}
	return this.task
	//return this.tasks[taskID]
}

func (this *RoleEventTaskManager) EnsureRefTask(refID int32) *db.DBRoleEventTask {
	//for _, task := range this.tasks {
	//	if task.RefID == refID {
	//		return task
	//	}
	//}
	if this.task == nil || this.task.RefID != refID {
		log.Error("time: %v userID: %v task :%v roleRefID: %v refID: %v",time.Now(),this.uid, this.task, this.task.RefID, refID)
		exception.GameException(exception.TASK_NOT_FOUND)
	}
	return this.task
}

//事件任务完成
func (this *RoleEventTaskManager) DoneEventTask(message *protocol.EventDone, starType int32) (*db.DBRoleEventTask, bool) {
	//result := []*db.DBRoleEventTask{}
	eventID := message.GetEventID()
	reward := message.GetReward()
	eventFaith := reward.GetFaith()
	eventBeliever := reward.GetBeliever()
	eventItemID := reward.GetItemID()

	var changeFaith int32 = 0
	believers := []*protocol.BelieverInfo{}
	if this.task == nil || this.task.RefID != eventID {
		return nil, false
	}

	task := this.task
	taskReward := conf.GetTaskReward(task.Type, task.EndingID)
	if taskReward != nil {
		task.RewardFaith = taskReward.FaithNum
		task.RewardBeliever = taskReward.BelieverInfo
		changeFaith += task.RewardFaith
		if task.RewardBeliever != nil {
			for _, believer := range task.RewardBeliever {
				believers = append(believers, believer)
			}
		}
	}

	task.RewardFaith = eventFaith
	if eventBeliever != nil {
		task.RewardBeliever = eventBeliever
	}
	if eventItemID != 0 {
		task.RewardItem = eventItemID
	}

	randomItem := eventItemID == 0 && taskReward.GetRelic != 0
	task.State = db.TASK_STATE_DONE
	this.task = nil
	return task, randomItem
}

//提交任务 返回奖励的信仰值
//func (this *RoleEventTaskManager) SubmitTask(taskID int32) int32 {
//	task := this.tasks[taskID]
//	if (task == nil) {
//		exception.GameException(exception.TASK_NOT_FOUND)
//	}
//	if (task.State != db.TASK_STATE_DONE) {
//		exception.GameException(exception.TASK_NOT_DONE)
//	}
//	delete(this.tasks, taskID)
//	//TODO 返回奖励id
//
//	return 1
//}

//获取一个不重复的任务id
func (this *RoleEventTaskManager) genTaskID() int32 {
	return 1
	//var maxID int32 = 0
	//for _, task := range this.tasks {
	//	if (task.ID > maxID) {
	//		maxID = task.ID
	//	}
	//}
	//return maxID + 1
}