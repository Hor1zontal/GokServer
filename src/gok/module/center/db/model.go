package db

import (
	"time"
	"gok/service/msg/protocol"
)

type DBNotice struct {
	ID     	int32 `bson:"_id" gorm:"AUTO_INCREMENT"` //公告id
	Title	string `bson:"title"`
	Content string `bson:"content"`
	Start	time.Time `bson:"start"`
	End		time.Time `bson:"end"`
	Status 	int32 `bson:"status" unique:"false"` // 0--进行中 1--未开始 2--已结束
}

func (notice *DBNotice) BuildProto() *protocol.Notice{
	return &protocol.Notice{
		NoticeID:notice.ID,
		Title:notice.Title,
		Content:notice.Content,
		Start:notice.Start.Unix(),
		End:notice.End.Unix(),
		Status:notice.Status,
	}
}
