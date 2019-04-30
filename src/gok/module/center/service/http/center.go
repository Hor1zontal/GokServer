package http

import (
	"aliens/common/character"
	"aliens/common/helper"
	"aliens/common/util"
	"github.com/name5566/leaf/log"
	"gok/constant"
	"gok/module/center/cache"
	"gok/module/center/conf"
	"gok/module/center/db"
	"gok/module/center/notice"
	"gok/service/msg/protocol"
	"gok/service/rpc"
	"net/http"
	"strconv"
	"time"
)

const (
	PARAM_METHOD = "method"
	PARAM_STATE = "state"
	PARAM_STATUS = "status"
	PARAM_MAINTAIN_STATE = "maintainState"
	PARAM_OPENTIME = "openTime"
	PARAM_ISCHECK_VERSION = "isCheckVersion"
	PARAM_QUERY_TYPE = "type"
	NOTICE_ID = "id"
	NOTICE_TITLE = "title"
	NOTICE_CONTENT = "content"
	NOTICE_START_TIME = "start"
	NOTICE_END_TIME = "end"
	NOTICE_STATUS = "status"
	WHITE_UID = "uid"
	PARAM_FLAG_KEY = "key"
	PARAM_FLAG_VALUE = "value"
)

func Center(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	method := request.FormValue(PARAM_METHOD)
	//sign := request.FormValue(PARAM_SIGN)
	response := ""

	//signText := method + conf.Server.AppKey
	succ := isSignSuccess(request.Form)
	if !succ {
		helper.GetResponse(RESULT_CODE_INVALID_SIGN)
		helper.SendToClient(responseWriter, response)
		return
	}
	switch method {
	//-------------------------中控接口
	case "change_state":
		response = changeServerState(request)
	case "get_state":
		response = getState(request)
	case "get_day_data":
		response = getDayData(request)
	case "get_user_data":
		response = getUserData(request)
	case "update_user_state":
		response = updateUserState(request)
	case "set_message_data":
		response = setMessageData(request)
	case "get_notice":
		response = getNotice(request)
	case "public_notice":
		response = publicNotice(request)
	case "delete_notice":
		response = deleteNotice(request)
	case "add_uid_white_list":
		response = addUidWhiteList(request)
	case "remove_uid_white_list":
		response = removeUidWhiteList(request)
	case "get_white_list_uid":
		response = getWhiteLisUid(request)
	case "get_star_summary":
		response = getStarSummary(request)
	case "refresh_client_version":
		response = refreshClientVersion(request)
	case "add_test_account":
		response = addTestAccount(request)
	case "clean_test_account":
		response = cleanTestAccount(request)
	case "get_test_account":
		response = getTestAccount()
	case "del_test_account":
		response = delTestAccount(request)
	case "update_star_flag":
		response = updateStarFlag(request)
	case "query_by_condition":
		response = queryByCondition(request)
	default:
		response = helper.GetResponse(RESULT_CODE_INVALID_PARAM)
	}

	helper.SendToClient(responseWriter, response)

}

func changeServerState(request *http.Request) string {
	maintainState := character.StringToInt32(request.FormValue(PARAM_MAINTAIN_STATE))
	if maintainState != constant.SERVER_MAINTAIN_STATE_OPEN && maintainState != constant.SERVER_MAINTAIN_STATE_CLOSE {
		return helper.GetResponse(RESULT_CODE_INVALID_PARAM)
	}
	serverState := character.StringToInt32(request.FormValue(PARAM_STATE))
	openTimeStamp := character.StringToInt64(request.FormValue(PARAM_OPENTIME))
	isCheckVersion,_ := strconv.ParseBool(request.FormValue(PARAM_ISCHECK_VERSION))
	//openTimestamp, err := time.ParseInLocation("2006-01-02 15:04:05", openTime, time.Local)
	//if err != nil {
	//	if openTime != "" {
	//		log.Debug("invalid server opentime : %v", err)
	//		return helper.GetResponse(RESULT_CODE_INVALID_PARAM)
	//	}
	//}
	if conf.ChangeState(serverState, openTimeStamp/1000, maintainState, isCheckVersion) {
		return helper.GetResponse(RESULT_CODE_SUCCESS)
	} else {
		return helper.GetResponse(RESULT_CODE_SERVER_EXCEPTION)
	}
}

func getState(request *http.Request) string {
	return helper.GetDataResponse(RESULT_CODE_SUCCESS, conf.Maintain)
}


type DayData struct {
	RegTotal	   int `json:"regTotal"`
	Login          int `json:"login"`
	Register       int `json:"register"`
	Charge         int `json:"charge"`
	ChargeFee      int `json:"chargeFee"`
}

func getDayData(request *http.Request) string {
	time := time.Now()
	data := DayData{}
	data.Login = cache.UserCache.GetDayLogin(time)
	data.Register = cache.UserCache.GetDayRegister(time)
	data.Charge = cache.UserCache.GetDayCharge(time)
	data.ChargeFee = cache.UserCache.GetDayChargeFee(time)
	data.RegTotal = cache.UserCache.GetRegisterTotal()
	return helper.GetDataResponse(RESULT_CODE_SUCCESS, data)
}

type UserData struct {
	UserID 		int32     `json:"userid" unique:"true"`  //用户id
	NickName    string    `json:"nickname" rorm:"nname"` //角色名称
	Channel     string    `json:"channel" rorm:"channel"`      //用户的渠道信息 渠道用户id存Username
	Power       int32     `json:"power" rorm:"power"`    //法力值
	PowerLimit  int32     `json:"powerLimit" rorm:"p_limit"`
	Faith       int32     `json:"faith" rorm:"faith"`    //信仰值   金币
	Diamond     int32     `json:"diamond" rorm:"diamond"`  //钻石
	Status  	byte      `json:"status" rorm:"status"`  //用户状态 0正常  1封号
}

func getUserData(request *http.Request) string {
	value := request.FormValue("value")
	//if uid == 0 {
	//	return helper.GetResponse(RESULT_CODE_INVALID_PARAM)
	//}
	queryType := character.StringToInt32(request.FormValue(PARAM_QUERY_TYPE))
	if queryType == 0 {
	//todo
	}
	var datas []*UserData
	switch queryType {
	case constant.QUERY_BY_UID:
		uid := character.StringToInt32(value)
		if !cache.UserCache.IsUserExist(uid) {
			return helper.GetResponse(RESULT_CODE_USER_NOTFOUND)
		}
		datas = getUserDataByUids([]int32{uid})
	case constant.QUERY_BY_NICKNAME:
		//rpc.UserServiceProxy.HandleMessage()
		//lpc.DBServiceProxy.QueryAllCondition(&db.DBNotice{}, "nickname", "", &result, db.DatabaseHandler)
		nickname := value
		uids := rpc.UserServiceProxy.QueryByNickname(nickname).GetUids()
		datas = getUserDataByUids(uids)
	case constant.QUERY_BY_USERNAME:
		username := value
		uids := rpc.PassportServiceProxy.QueryByUsername(username).GetUids()
		datas = getUserDataByUids(uids)
	}
	return helper.GetDataResponse(RESULT_CODE_SUCCESS, datas)
}

func getUserDataByUids(uids []int32) []*UserData {
	var datas []*UserData
	for _, uid := range uids {
		data := &UserData{UserID:uid}
		cache.UserCache.HGetUser(uid, data)
		datas = append(datas, data)
	}
	return datas
}

func updateUserState(request *http.Request) string {
	//memo := request.FormValue(PARAM_MEMO)
	uidStr := request.FormValue(PARAM_UID)
	statusStr := request.FormValue(PARAM_STATUS)
	//sign := request.FormValue(PARAM_SIGN)
	uid := character.StringToInt32(uidStr)
	if uid == 0 {
		return helper.GetResponse(RESULT_CODE_INVALID_PARAM)
	}
	status := character.StringToInt32(statusStr)

	//signText := "uid" + uidStr + "&status=" + statusStr + conf.Server.AppKey

	//signResult := cipher.MD5Hash(signText)
	//
	//if signResult != sign {
	//	log.Debug("invalid sign ! signResult %v : sign : %v", signResult, sign)
	//	//测试跳过签名
	//	//return RESULT_CODE_INVALID_SIGN
	//}
	if !cache.UserCache.IsUserExist(uid) {
		return helper.GetResponse(RESULT_CODE_USER_NOTFOUND)
	}
	rpc.PassportServiceProxy.UserState(uid, status)
	return helper.GetResponse(RESULT_CODE_SUCCESS)
}

func setMessageData(request *http.Request) string {
	messageData := request.FormValue("message")
	message := &protocol.LampMessage{
		Data: messageData,
	}
	rpc.UserServiceProxy.BroadcastAll(&protocol.GS2C{Sequence:[]int32{1071},LampMessagePush:message})
	return helper.GetResponse(RESULT_CODE_SUCCESS)
}

func getNotice(request *http.Request) string {
	noticeStatus := character.StringToInt32(request.FormValue("status"))

	var result []*db.DBNotice
	db.DatabaseHandler.QueryAllCondition(&db.DBNotice{}, "status", noticeStatus, &result)
	if result == nil {
		return helper.GetResponse(RESULT_CODE_DATABASE_EXCEPTION)
	}
	return helper.GetDataResponse(RESULT_CODE_SUCCESS, result)
}

func publicNotice(request *http.Request) string {
	title := request.FormValue(NOTICE_TITLE)
	content := request.FormValue(NOTICE_CONTENT)
	//start := request.FormValue(NOTICE_START_TIME)
	//end := request.FormValue(NOTICE_END_TIME)
	status := character.StringToInt32(request.FormValue(NOTICE_STATUS))
	start := character.StringToInt64(request.FormValue(NOTICE_START_TIME))/1000
	end := character.StringToInt64(request.FormValue(NOTICE_END_TIME))/1000

	startTime := util.GetTime(start)
	endTime := util.GetTime(end)

	noticeData := &db.DBNotice{
		Title:title,
		Content:content,
		Start:startTime,
		End:endTime,
		Status:status,
	}
	err := db.DatabaseHandler.Insert(noticeData)
	if err != nil {
		log.Debug(err.Error())
	}

	if status == constant.NOTICE_WILL {
		notice.NoticesManager.WillNotices[noticeData.ID] = noticeData
	}
	return helper.GetResponse(RESULT_CODE_SUCCESS)
}

func deleteNotice(request *http.Request) string {
	noticeID := request.FormValue(NOTICE_ID)
	notice.NoticesManager.DeleteNotice(character.StringToInt32(noticeID))
	err := db.DatabaseHandler.DeleteOne(&db.DBNotice{ID:character.StringToInt32(noticeID)})
	if err != nil {
		return helper.GetDataResponse(RESULT_CODE_DATABASE_EXCEPTION, err.Error())
	}
	return helper.GetResponse(RESULT_CODE_SUCCESS)
}

func addUidWhiteList(request *http.Request) string {
	uid := character.StringToInt32(request.FormValue(WHITE_UID))
	if uid <= 0 {
		return helper.GetResponse(RESULT_CODE_INVALID_PARAM)
	}
	if !cache.UserCache.IsUserExist(uid) {
		return helper.GetResponse(RESULT_CODE_USER_NOTFOUND)
	}
	if cache.ClusterCache.IsUidWhiteList(uid) {
		return helper.GetResponse(RESULT_CODE_USER_EXIST)
	}
	if !cache.ClusterCache.AddUidWhiteList(uid) {
		return helper.GetResponse(RESULT_CODE_SERVER_EXCEPTION)
	}
	return helper.GetResponse(RESULT_CODE_SUCCESS)
}

func removeUidWhiteList(request *http.Request) string {
	uid := character.StringToInt32(request.FormValue(WHITE_UID))
	if uid <= 0 {
		return helper.GetResponse(RESULT_CODE_INVALID_PARAM)
	}
	if !cache.ClusterCache.RemoveUidWhiteList(uid) {
		return helper.GetResponse(RESULT_CODE_SERVER_EXCEPTION)
	}
	return helper.GetResponse(RESULT_CODE_SUCCESS)
}

func getWhiteLisUid(request *http.Request) string {
	uids := cache.ClusterCache.GetWhiteListUid()
	return helper.GetDataResponse(RESULT_CODE_SUCCESS, uids)
}


func getStarSummary(request *http.Request) string {
	data := cache.ClusterCache.GetStarSummary()
	return helper.GetDataResponse(RESULT_CODE_SUCCESS, data)
}

func refreshClientVersion(request *http.Request) string {
	resp := rpc.PassportServiceProxy.RefreshClientVersion()
	if resp == nil {
		return helper.GetResponse(RESULT_CODE_SERVER_EXCEPTION)
	}
	if !resp.Result {
		return helper.GetResponse(RESULT_CODE_FAILED)
	}
	return helper.GetResponse(RESULT_CODE_SUCCESS)

}

func addTestAccount(request *http.Request) string {
	uid := character.StringToInt32(request.FormValue(WHITE_UID))
	if uid <= 0 {
		return helper.GetResponse(RESULT_CODE_INVALID_PARAM)
	}
	if !cache.UserCache.IsUserExist(uid) {
		return helper.GetResponse(RESULT_CODE_USER_NOTFOUND)
	}
	if cache.UserCache.ExistTestUserID(uid) {
		return helper.GetResponse(RESULT_CODE_USER_EXIST)
	}
	if !cache.UserCache.SetTestUserID(uid) {
		return helper.GetResponse(RESULT_CODE_FAILED)
	}
	return helper.GetResponse(RESULT_CODE_SUCCESS)
}

func cleanTestAccount(request *http.Request) string {
	uid := request.FormValue(PARAM_UID)
	resp := rpc.PassportServiceProxy.CleanTestAccount(character.StringToInt32(uid))
	if resp == nil {
		return helper.GetResponse(RESULT_CODE_SERVER_EXCEPTION)
	}
	if !resp.GetResult() {
		return helper.GetDataResponse(RESULT_CODE_FAILED, resp.GetMessage())
	}
	return helper.GetResponse(RESULT_CODE_SUCCESS)
}

func getTestAccount() string {
	return helper.GetDataResponse(RESULT_CODE_SUCCESS, cache.UserCache.GetTestUserID())
}

func delTestAccount(request *http.Request) string {
	uid := character.StringToInt32(request.FormValue(PARAM_UID))
	if !cache.UserCache.DelTestUserID(uid) {
		return helper.GetResponse(RESULT_CODE_FAILED)
	}
	return helper.GetResponse(RESULT_CODE_SUCCESS)
}

func updateStarFlag(request *http.Request) string {
	key := character.StringToInt32(request.FormValue(PARAM_FLAG_KEY))
	value := character.StringToInt32(request.FormValue(PARAM_FLAG_VALUE))
	result := rpc.StarServiceProxy.UpdateAllStarFlag(key, value).GetResult()
	return helper.GetDataResponse(RESULT_CODE_SUCCESS, result)
}

func queryByCondition(request *http.Request) string {
	start := character.StringToInt64(request.FormValue("start"))/1000
	end := character.StringToInt64(request.FormValue("end"))/1000
	buildLv := character.StringToInt32(request.FormValue("buildLv"))
	limit := character.StringToInt32(request.FormValue("limit"))
	skip := character.StringToInt32(request.FormValue("skip"))

	datas := rpc.StarServiceProxy.GetUidsByCondition(start, end, buildLv, limit, skip)
	starDatas := datas.GetUserData()
	count := datas.GetCount()
	users := make([]*UserDetail, len(starDatas))
	for index, starData := range starDatas {
		data := &UserData{UserID:starData.Uid}
		cache.UserCache.HGetUser(starData.Uid, data)
		users[index] = &UserDetail{
			UserData:data,
			StarData:starData,
		}
	}
	return helper.GetDataResponse(RESULT_CODE_SUCCESS, &UsersDetailData{Count:count,Detail:users})
}

type UsersDetailData struct {
	Count 		int32 			`json:"count"`
	Detail      []*UserDetail	`json:"users"`
}

type UserDetail struct {
	UserData	*UserData					`json:"userData"`
	StarData	*protocol.UserStarData		`json:"starData"`

}