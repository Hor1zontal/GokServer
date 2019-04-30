package manager

import (
	"gok/module/game/db"
	//"gok/service/exception"
	//"gok/service/msg/protocol"
)

//角色标识
type RoleStarManager struct {
	uid      int32
	starID   int32
	starType int32
	starSize int32
	//starsRecordOri      []int32 //用户原始星球收藏 wjl 20170602
	//starsRecordUser     []int32 //用户占用星球收藏 wjl 20170602
	//starsRecordComplete []int32 //用户已经完成的星球数据 wjl 20170602
}

//初始化
func (this *RoleStarManager) Init(role *db.DBRole) {
	this.uid = role.UserID

	//this.starsRecordOri = role.StarsRecordOri
	//this.starsRecordUser = role.StarsRecordUser
	//this.starsRecordComplete = role.StarsComplete
}

//更新数据库内存
func (this *RoleStarManager) Update(role *db.DBRole) {
	//role.Stars = this.stars
	//role.StarsRecordOri = this.starsRecordOri
	//role.StarsRecordUser = this.starsRecordUser
	//role.StarsComplete = this.starsRecordComplete
}
//
////获取原始星球记录列表 wjl 20170602
//func (this *RoleStarManager) GetRecordOri() []int32 { return this.starsRecordOri }
//
////获取占领星球记录列表 wjl 20170602
//func (this *RoleStarManager) GetRecordUser() []int32 { return this.starsRecordUser }

//获取当前拥有的星球id
func (this *RoleStarManager) GetStarId() int32 {
	return this.starID
}

func (this *RoleStarManager) GetStarType() int32 {
	return this.starType
}

func (this *RoleStarManager) UpdateStarInfo(starID int32, starType int32) {
	this.starID = starID
	this.starType = starType
	//this.starSize = starSize
}

//func (this *RoleStarManager) GetStarSize() int32 {
//	return this.starSize
//}

//func (this *RoleStarManager)(info protocol.LoginStarInfoRet) {
//	if info.GetCurrStar() != nil {
//		this.starType = info.GetCurrStar().GetStarID()
//		this.starID = info.GetCurrStar().GetStarType()
//	}
//	for _, star := range info.Star {
//		this.starHistory
//	}
//}

////新增星球
//func (this *RoleStarManager) AddStar(starId int32, starType int32) {
//	if (this.ContainsStar(starId)) {
//		exception.GameException(exception.STAR_ALREADY_OWN)
//	}
//	if this.stars == nil{
//		this.stars = &db.DBRoleStar{};
//	}
//	this.stars.ID = starId;
//	this.stars.Type = starType;
////	this.stars = append(this.stars, &db.DBRoleStar{ID:starId, Type:starType})
//}

//目前是否拥有开发中的星球
//func (this *RoleStarManager) ContainsStar(starId int32) bool {
//	if this.stars == nil{ return false };
//	if this.stars.ID == starId {
//		return true
//	}
//	return false;
//	//for _, star := range this.stars {
//	//	if (star.ID == starId) {
//	//		return true
//	//	}
//	//}
////	return false
//}

//func (this *RoleStarManager) GetStarRecord() { //获取星球记录
//
//}

//func (this *RoleStarManager) IsOriStarRecordMax() bool { //原始星球记录是否达到最大值
//	return len(this.starsRecordOri) >= 5
//}
//
//func (this *RoleStarManager) IsUserStarRecordMax() bool { //用户星球记录是否达到最大值
//	return len(this.starsRecordUser) >= 20
//}
//
//func (this *RoleStarManager) ContainsOriStarRecord(recordID int32) bool { //是否包含原始星球记录
//	for _, v := range this.starsRecordOri {
//		if v == recordID {
//			return true;
//		}
//	}
//	return false
//}
//
//func (this *RoleStarManager) ContainsUserStarRecord(recordID int32) bool { //是否包含已占领星球记录
//	for _, v := range this.starsRecordUser {
//		if v == recordID {
//			return true;
//		}
//	}
//	return false
//}

//func (this *RoleStarManager) UpdateStarRecord(recordType int32, recordID int32) { //更新星球记录 wjl 20170602
//
//	if recordType == 0x00 { //原始星球
//
//		if (this.IsOriStarRecordMax()) {
//			exception.GameException(exception.STAR_RECORD_FAILED_MAX)
//		}
//
//		if (this.ContainsOriStarRecord(recordID)) {
//			exception.GameException(exception.STAR_RECORD_FAILED_EXIST)
//		}
//		this.starsRecordOri = append(this.starsRecordOri, recordID);
//	}
//
//	if recordType == 0x01 { //已被占领的星球
//		if (this.IsUserStarRecordMax()) {
//			exception.GameException(exception.STAR_RECORD_FAILED_MAX)
//		}
//
//		if (this.ContainsUserStarRecord(recordID)) {
//			exception.GameException(exception.STAR_RECORD_FAILED_EXIST)
//		}
//		this.starsRecordUser = append(this.starsRecordUser, recordID);
//	}
//}
//
//func (this *RoleStarManager) DelStarRecord(recordType int32, recordID int32) bool { //请求删除星球记录 wjl 20170605
//	if recordType == 0x00 { //原始星球
//		for i, v := range this.starsRecordOri {
//			if v == recordID {
//				this.starsRecordOri = append(this.starsRecordOri[:i], this.starsRecordOri[i+1:]...)
//				return true;
//			}
//		}
//		return false;
//	}
//
//	if recordType == 0x01 { //已被占领的星球
//		for i, v := range this.starsRecordUser {
//			if v == recordID {
//				this.starsRecordUser = append(this.starsRecordUser[:i], this.starsRecordUser[i+1:]...)
//				return true;
//			}
//		}
//		return false;
//	}
//	return false;
//}

//替换星球记录信息
//func (this *RoleStarManager) ReplaceStarRecord(recordType int32, recordID int32, replaceID int32) {
//	if recordType == 0x00 { //原始星球
//		if (this.ContainsOriStarRecord(replaceID)) {
//			exception.GameException(exception.STAR_RECORD_FAILED_EXIST)
//		}
//		for i, v := range this.starsRecordOri {
//			if v == recordID {
//				this.starsRecordOri[i] = replaceID
//				return
//			}
//		}
//	} else if recordType == 0x01 { //已被占领的星球
//		if (this.ContainsUserStarRecord(replaceID)) {
//			exception.GameException(exception.STAR_RECORD_FAILED_EXIST)
//		}
//		for i, v := range this.starsRecordUser {
//			if v == recordID {
//				this.starsRecordUser[i] = replaceID
//				return
//			}
//		}
//	}
//	exception.GameException(exception.STAR_RECORD_FAILED_NOTFOUND)
//}

//func (this *RoleStarManager) MoveStarRecord(srcType int32, srcID int32, destType int32, destID int32) bool { //wjl 20170607 移动星球记录
//	if srcType == destType { //同类型列表不允许 移动
//		return false;
//	}
//
//	if srcType == 0x00 { //从原始星球列表 移动到  用户星球列表
//
//		for _, v := range this.starsRecordUser {
//			if v == destID { //表示是存在的 恩恩 那这里问题很严重呀
//				return false;
//			}
//		}
//
//		for i, v := range this.starsRecordOri {
//			if v == srcID { //表示是存在的 恩恩 这就没什么问题了
//				this.starsRecordOri = append(this.starsRecordOri[:i], this.starsRecordOri[i+1:]...); //列表中移除
//				break;
//			}
//		}
//		this.starsRecordUser = append(this.starsRecordUser, destID); //插入列表内
//		return true;
//	}
//
//	if srcType == 0x01 { //从用户星球列表移动到原始星球列表( 理论上不允许这个操作 )
//		return false;
//	}
//	return false;
//}

//func (this *RoleStarManager) Record(starType int32) {
//	if starType == 0 {
//		return
//	}
//	this.starsRecordComplete = append(this.starsRecordComplete, starType)
//}

//func (this *RoleStarManager) OccupyStar(srcStar *protocol.StarInfoDetail, destStar *protocol.StarInfoDetail) {
//	this.Record(srcStar.GetType())
//	//dbStarBuilding := []*db.DBBuilding{};
//
//	//for _, v := range srcStar.Building{
//	//	dbStarBuilding = append( dbStarBuilding, &db.DBBuilding{
//	//		ID: v.GetId(),
//	//		Type:v.GetType(),
//	//		Level:v.GetLevel(),
//	//	})
//	//}
//
//	//dbStarComplete := &db.DBRoleStar_Complete{
//	//	ID: srcStar.GetStarID(),
//	//	Type: srcStar.GetType(),
//	//	Building:dbStarBuilding,
//	//}
//	//this.starsRecordComplete = append( this.starsRecordComplete, dbStarComplete );
//	//this.stars.ID = destStar.GetStarID();
//	//this.stars.Type = destStar.GetType();
//
//	for i, v := range this.starsRecordOri { //判断是否在原始星球中是否存在 如果存在 则从原始列表中删除
//		if v == destStar.GetStarID() {
//			this.starsRecordOri = append(this.starsRecordOri[:i], this.starsRecordOri[i+1:]...); //列表中移除
//			break;
//		}
//	}
//}
