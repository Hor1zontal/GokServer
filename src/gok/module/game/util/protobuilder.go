package util

import (
	"gok/service/msg/protocol"
	"gok/constant"
	"time"
	"gok/module/game/conf"
	"gok/module/game/db"
)

func IsOverdue(newsFeed *protocol.NewsFeed) bool {
	return time.Now().Unix() - newsFeed.GetTime() > conf.DATA.CountdownApplying
}

func BuildBuildingInfoPush(uid int32, starID int32, starType int32, buildings []*protocol.BuildingInfo) *protocol.GS2C {
	message := &protocol.BuildingInfoPush{
		Uid:  uid,//拥有者的用户ID
		StarID:  starID,
		Type:  starType,
		Building:buildings,
	}
	return &protocol.GS2C{
		Sequence:[]int32{1031},
		BuildingInfoPush:message,
	}
}

//构建踢出用户消息
func BuildKickPush(kickType constant.LOGOUT_TYPE) *protocol.GS2C {
	return &protocol.GS2C{
		Sequence:[]int32{1001},
		KickoffPush:&protocol.KickoffPush{
			Type:int32(kickType),
		}}
}


//构建角色信息
//func BuildDBRoleInfo(this *db.DBRole) *protocol.RoleInfo {
//	return &protocol.RoleInfo{
//		Id:this.UserID),
//		Icon:this.Icon),
//		Nickname:this.NickName),
//		Level:this.Level),
//		Exp:this.Exp),
//	}
//}


func BuildOnlinePush(uid int32) *protocol.GS2C {
	return &protocol.GS2C{Sequence:[]int32{1055},OnlinePush: uid}
}

func BuildOfflinePush(uid int32) *protocol.GS2C {
	return &protocol.GS2C{Sequence:[]int32{1056}, OfflinePush: uid}
}

func BuildItemPush(item *db.DBRoleItem) *protocol.GS2C {
	return &protocol.GS2C{
		Sequence:[]int32{1050},
		ItemPush: &protocol.BagItem{
			Id:  item.ID,
			Num: item.Num,
		}}
}


func BuildTempItemPush(itemID int32) *protocol.GS2C {
	return &protocol.GS2C{
		Sequence:[]int32{1058},
		TempItemPush: itemID}
}

//构建物品出售通知
//func BuildSalePush(sale *protocol.Sale) *protocol.GS2C {
//	return &protocol.GS2C{
//		Session: 1054),
//		SalePush: sale}
//}

func BuildStarFlagPush(flag *protocol.FlagInfo) *protocol.GS2C {
	message := flag
	return &protocol.GS2C{
		Sequence:[]int32{1082},
		StarFlagPush:message,
	}
}