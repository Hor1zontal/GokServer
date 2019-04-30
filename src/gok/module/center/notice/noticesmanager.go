package notice

import (
	"time"
	"gok/module/center/db"
	"gok/constant"
	"gok/service/rpc"
	"gok/service/msg/protocol"
)

var NoticesManager = &Notices{
	OnNotices:make(map[int32]*db.DBNotice),
	WillNotices:make(map[int32]*db.DBNotice),
	UpdateNotices:make(map[int32]*db.DBNotice),
}

type Notices struct {
	OnNotices map[int32]*db.DBNotice
	WillNotices map[int32]*db.DBNotice
	UpdateNotices map[int32]*db.DBNotice
}

func Init() {
	NoticesManager.LoadData()
}

func (this *Notices)LoadData(){
	var notices []*db.DBNotice
	db.DatabaseHandler.QueryAll(&db.DBNotice{}, &notices)
	for _, notice := range notices {
		if notice.Status == constant.NOTICE_ON {
			this.OnNotices[notice.ID] = notice
		}
		if notice.Status == constant.NOTICE_WILL {
			this.WillNotices[notice.ID] = notice
		}
	}
}

func (this *Notices)updateData(){

	for _,notice := range this.UpdateNotices {
		db.DatabaseHandler.UpdateOne(notice)
		delete(this.UpdateNotices, notice.ID)
	}

}

func (this *Notices)DealTimeout() {
	now := time.Now()
	if this.WillNotices != nil || len(this.WillNotices) != 0 {
		this.dealWillNotices(now)
	}
	if this.OnNotices != nil || len(this.OnNotices) != 0 {
		this.dealOnNotices(now)
	}
	if this.UpdateNotices != nil || len(this.UpdateNotices) != 0 {
		this.updateData()
	}
}

func (this *Notices)updateNoticeStatus(notice *db.DBNotice, status int32) {
	notice.Status = status
	if status == constant.NOTICE_ON {
		delete(this.WillNotices, notice.ID)
		this.OnNotices[notice.ID] = notice
	}
	if status == constant.NOTICE_DONE {
		delete(this.OnNotices, notice.ID)
	}
	this.UpdateNotices[notice.ID] = notice
}

func (this *Notices)dealWillNotices(now time.Time) {
	for _, notice := range this.WillNotices {
		ret := notice.Start.Local().After(now)
		if !ret{
			this.updateNoticeStatus(notice, constant.NOTICE_ON)
			rpc.UserServiceProxy.BroadcastAll(&protocol.GS2C{Sequence:[]int32{1072},NoticeMessagePush:&protocol.NoticeMessage{NoticeID:notice.ID}})
		}
	}
}

func (this *Notices)dealOnNotices(now time.Time) {
	for _, notice := range this.OnNotices {
		if notice.End.Local().Before(now) {
			this.updateNoticeStatus(notice, constant.NOTICE_DONE)
		}
	}
}

func (this *Notices)DeleteNotice(noticeID int32) {
	_, ok := this.OnNotices[noticeID]
	if ok {
		delete(this.OnNotices, noticeID)
	}
	_, ok1 := this.WillNotices[noticeID]
	if ok1 {
		delete(this.WillNotices, noticeID)
	}
	_, ok2 := this.UpdateNotices[noticeID]
	if ok2 {
		delete(this.UpdateNotices, noticeID)
	}
}