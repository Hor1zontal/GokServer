package cache

import (
	"aliens/common/character"
	"gok/service/msg/protocol"
	"aliens/log"
)

const (
	//MAILS_USER_RECEIVE_PREFIX = "mur_"	//存储用户的邮件


	USER_MAILS_KEY = "userMail:" //用户拥有的邮件的key
	MAIL_KEY = "mail:" //邮件的key

	//FLAG_LOAD_USERMAILS string = "flumails_" //标识，是否加载用户所拥有的邮件到内存
	FLAG_LOAD_MAIL string = "flag:mail"  	 	    //标识，是否加载邮件信息到内存
)


type MailCacheManager struct {
	*cacheManager
}

func NewMailCacheManager() *MailCacheManager {
	return &MailCacheManager{
		&cacheManager{},
	}
}

func getUserMailKey(id int32) string {
	return USER_MAILS_KEY + character.Int32ToString(id)
}

func getMailKey(id int64) string {
	return MAIL_KEY + character.Int64ToString(id)
}

//func (this *MailCacheManager) GetUserMail(userID int32, mailID int64) bool{
//	return this.redisClient.SDelData(getUserMailKey(userID), mailID)
//}



func (this *MailCacheManager) SetMailCache(mailID int64, mailInfo *protocol.Mail) bool {
	bytes, err := mailInfo.Marshal()
	if err != nil {
		log.Debug("mail bson marshal error: %v", err.Error())
		return false
	}
	return this.redisClient.SetData(getMailKey(mailID), bytes)
}

func (this *MailCacheManager) GetMailCache(mailID int64) *protocol.Mail{
	bytes := this.redisClient.GetBytesData(getMailKey(mailID))
	if bytes == nil || len(bytes) == 0{
		return nil
	}
	result := &protocol.Mail{}
	err := result.Unmarshal(bytes)
	if err != nil {
		log.Debug("mail bson unmarshal error: %v", err.Error())
		return nil
	}
	return result
}

func (this *MailCacheManager) DelMailCache(mailID int64) bool {
	return this.redisClient.DelData(getMailKey(mailID))
}

func (this *MailCacheManager) SetExpireTime(mailID int64, seconds int) bool {
	return this.redisClient.Expire(getMailKey(mailID), seconds)
}

//TODO 用SortedSet 支持分页
//func(this *MailCacheManager) GetUserMails(uid int32, limitTime int32, offset int32, count int32) []string{
//	return this.redisClient.ZRevRangeByScoreBeforeLimit(getUserMailKey(uid), limitTime, offset, count)
//}
//
//func(this *MailCacheManager) GetUserMail(uid int32, mailID string) bool {
//	return this.redisClient.ZRem(getUserMailKey(uid), mailID)
//}
//
//func(this *MailCacheManager) SetUserMail(uid int32, mailID string, createTimestamp time.Time) bool {
//	return this.redisClient.ZAdd(getUserMailKey(uid), createTimestamp.Unix(), mailID)
//}
//
//func(this *MailCacheManager) RemoveUserMail(uid int32, mailID string) bool {
//	return this.redisClient.ZRem(getUserMailKey(uid), mailID)
//}

func (this *MailCacheManager) GetUserMails(userID int32) []int{
	return this.redisClient.SMembers(getUserMailKey(userID))
}

func (this *MailCacheManager) GetUserMail(userID int32, mailID int64) bool {
	return this.redisClient.SContains(getUserMailKey(userID), mailID)
}

func (this *MailCacheManager) SetUserMail(userID int32, mailID int64) bool {
	return this.redisClient.SAddData(getUserMailKey(userID), mailID)
}

func (this *MailCacheManager) DelUserMail(userID int32, maiID int64) bool {
	return this.redisClient.SDelData(getUserMailKey(userID), maiID)
}