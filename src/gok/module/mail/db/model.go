package db

import (
	"time"
	"gok/service/msg/protocol"
)

const (
//collection constants
//increment id constants
//ID_USER = "uid
)

type DBMail struct {
	ID          int64 	`bson:"_id" gorm:"PRIMARY_KEY;AUTO_INCREMENT"`    	 //邮件id
	Owner		int32 	`bson:"owner"`  //邮件所属用户 0代表全局邮件
	Title       string  `bson:"title"`  //邮件标题
	Content		string	`bson:"content"`  //邮件内容
	CreateTime  time.Time 	   `bson:"create_time"` //邮件创建时间
	Attach		*DBMailAttach  `bson:"attach"`     		//邮件附件
}

func (this *DBMail) BuildProtocol() *protocol.Mail {
	mail := &protocol.Mail{}
	mail.Id = this.ID
	mail.Owner = this.Owner
	mail.Title = this.Title
	mail.Content = this.Content
	mail.CreateTime = this.CreateTime.Unix()
	if this.Attach != nil {
		mail.Attach = this.Attach.BuildProtocol()
	}
	return mail
}

//邮件附件
type DBMailAttach struct {
	Draw	  bool			`bson:"draw"`     		//是否领取
	Power     int32         `bson:"power"`     		//法力值
	Faith     int32         `bson:"faith"`     		//信仰值
	GayPoint  int32         `bson:"gay_point"` 	    //友情点
	Diamond	  int32			`bson:"diamond"`     	//钻石
	Believer  []*Believer   `bson:"believer"`		//信徒
	Item      []*Item   	`bson:"item"`			//物品 图鉴
}

func (this *DBMailAttach) BuildProtocol() *protocol.Attach {
	attach := &protocol.Attach{}
	attach.Draw = this.Draw
	attach.Power = this.Power
	attach.Faith = this.Faith
	attach.GayPoint = this.GayPoint
	attach.Diamond = this.Diamond

	if this.Believer != nil {
		believers := []*protocol.BelieverInfo{}
		for _, believer := range this.Believer {
			believers = append(believers, &protocol.BelieverInfo{Id:believer.ID, Num:believer.Num})
		}
		attach.Believer = believers
	}

	if this.Item != nil {
		items := []*protocol.BagItem{}
		for _, item := range this.Item {
			items = append(items, &protocol.BagItem{Id:item.ID, Num:item.Num})
		}
		attach.Item = items
	}

	return attach
}

//信徒
type Believer struct {
	ID 		string   `bson:"_id" json:"id"`     		//信徒id
	Num 	int32	 `bson:"num" json:"num"`     		//信徒数量
}

//圣物
type Item struct {
	ID 		int32   `bson:"_id" json:"id"`     		//信徒id
	Num 	int32	 `bson:"num" json:"num"`     		//信徒数量
}
