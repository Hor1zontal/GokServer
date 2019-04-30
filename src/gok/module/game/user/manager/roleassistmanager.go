package manager

//import (
//	"gok/module/game/db"
//	"time"
//)
//
////用户协助请求管理
//type RoleAssistManager struct {
//	assists map[int32]*db.DBRoleAssist   //事件id和事件招募请求的对应关系
//}
//
////初始化
//func (this *RoleAssistManager) Init(role *db.DBRole) {
//	this.assists = make(map[int32]*db.DBRoleAssist)
//	for _, assist := range role.Assist {
//		this.assists[assist.EventID] = assist
//	}
//}
//
////更新内存
//func (this *RoleAssistManager) Opt(role *db.DBRole) {
//	role.Assist = this.GetAssists()
//}
//
//func (this *RoleAssistManager) GetAssists() []*db.DBRoleAssist {
//	result := []*db.DBRoleAssist{}
//	for _, assist := range this.assists {
//		result = append(result, assist)
//	}
//	return result
//}
//
////新建用户事件招募
//func (this *RoleAssistManager) NewAssist(eventID int32, uid int32, nickname string, msg string) *db.DBRoleAssist {
//	assist := this.assists[eventID]
//	if (assist == nil) {
//		assist = &db.DBRoleAssist{
//			EventID:eventID,
//			Uid:uid,
//			NickName:nickname,
//			Msg:msg,
//			CreateTime:time.Now(),
//		}
//		this.assists[eventID] = assist
//	}
//	return assist
//}
//
////删除事件协助
//func (this *RoleAssistManager) DeleteAssist(eventID int32) {
//	delete(this.assists, eventID)
//}