package manager
//
//import (
//	"gok/module/game/db"
//	"gok/service/msg/protocol"
//	"github.com/gogo/protobuf/proto"
//	"aliens/common/util"
//	"gok/constant"
//	"gok/module/game/cache"
//)
//
//
////角色标识
//type RoleStrangerManager struct {
//	strangers map[int32]*db.DBStranger
//}
//
////初始化
//func (this *RoleStrangerManager) Init(role *db.DBRole) {
//	this.strangers = make(map[int32]*db.DBStranger)
//	for _, stranger := range role.Strangers {
//		this.strangers[stranger.ID] = stranger
//	}
//}
//
////更新数据库内存
//func (this *RoleStrangerManager) Opt(role *db.DBRole) {
//	role.Strangers = this.GetStrangers()
//}
//
////获取所有物品
//func (this *RoleStrangerManager) GetStrangers() []*db.DBStranger {
//	result := []*db.DBStranger{}
//	for _, stranger := range this.strangers {
//		result = append(result, stranger)
//	}
//	return result
//}
//
//func (this *RoleStrangerManager) RemoveLatestActiveStranger() {
//	var lastStranger *db.DBStranger = nil
//	for _, stranger := range this.strangers {
//		this.strangers[stranger.ID] = stranger
//		if lastStranger == nil || lastStranger.ActiveTime.After(stranger.ActiveTime) {
//			lastStranger = stranger
//		}
//	}
//	if (lastStranger != nil) {
//		delete(this.strangers, lastStranger.ID)
//	}
//}
//
//func (this *RoleStrangerManager) AddStranger(stranger *protocol.Stranger) {
//	if (len(this.strangers) >= constant.STRANGER_LIMIT) {
//		this.RemoveLatestActiveStranger()
//	}
//	this.strangers[stranger.GetId()] = &db.DBStranger{
//		ID:stranger.GetId(),
//		Type:stranger.GetType(),
//		Param:stranger.GetParam(),
//		ActiveTime:util.GetTime(stranger.GetTime()),
//	}
//}
//
//func (this *RoleStrangerManager) RemoveStranger(id int32) {
//	delete(this.strangers, id)
//}
//
////获取所有物品
//func (this *RoleStrangerManager) GetProtocolStrangers() []*protocol.Stranger {
//	result := []*protocol.Stranger{}
//	for _, stranger := range this.strangers {
//		result = append(result, &protocol.Stranger{
//			Id:stranger.ID),
//			Type:stranger.Type),
//			Param:stranger.Param),
//			Nickname:cache.UserCache.GetUserNickname(stranger.ID)),
//			Time:stranger.ActiveTime.Unix()),
//			Avatar:cache.UserCache.GetUserAvatar(stranger.ID)),
//		})
//	}
//	return result
//}
