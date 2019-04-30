/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/4/17
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package mail

import (
	"gok/module/mail/db"
	"time"
	"gok/service/msg/protocol"
	"gopkg.in/mgo.v2/bson"
	"aliens/log"
	"gok/service/lpc"
	"gok/module/mail/cache"
	"gok/service/exception"
	"gok/module/mail/conf"
)

const (
	MAIL_LIMIT = 20
	RESULT_CODE_INVALID_PARAM = 1
)

//var Manager = &_mailManager{userMailMapping:make(map[int32]map[int64]*protocol.Mail)}
var Manager = &mailManager{}

type mailManager struct {
	//userMailMapping map[int32]map[int64]*protocol.Mail
}

func (this *mailManager) CreateMail(owner int32, title string, content string, attach string) *protocol.Mail {
	var mailAttach *db.DBMailAttach = nil
	if attach != "" {
		mailAttach = &db.DBMailAttach{}
		error := bson.UnmarshalJSON([]byte(attach), mailAttach)
		if error != nil {
			log.Debug("attach parse error : %v", error)
			return nil
		}
	}
	mail := &db.DBMail{
		Owner:owner,
		Title:title,
		Content:content,
		CreateTime:time.Now(),
		Attach:mailAttach,
	}
	err1 := db.DatabaseHandler.Insert(mail)
	if err1 != nil {
		return nil
	}
	mailInfo := mail.BuildProtocol()

	cache.MailCache.SetUserMail(mail.Owner, mail.ID)
	cache.MailCache.SetMailCache(mail.ID, mailInfo)
	cache.MailCache.SetExpireTime(mail.ID, conf.Server.RedisExpireTime)
	return mailInfo
}

func cleanMail(mails map[int64]*protocol.Mail) {
	var minMail int64 = 0
	for id, _ := range mails {
		if minMail == 0 {
			minMail = id
		} else if id < minMail {
			minMail = id
		}
	}
	if minMail != 0 {
		delete(mails, minMail)
		lpc.DBServiceProxy.Delete(&db.DBMail{ID:minMail}, db.DatabaseHandler)
	}
}

func (this *mailManager) EnsureUserMail(owner int32, mailID int64) {
	if !cache.MailCache.GetUserMail(owner, mailID) {
		exception.GameException(exception.MAIL_NOT_FOUND)
	}
}

func (this *mailManager) DrawMail(owner int32, mailID int64) *protocol.Mail {

	this.EnsureUserMail(owner, mailID)
	mail := cache.MailCache.GetMailCache(mailID)
	if mail == nil {
		//exception.GameException(exception.MAIL_NOT_FOUND)
		cache.MailCache.DelUserMail(owner, mailID)
	}
	return mail
}

func (this *mailManager) RemoveMail(owner int32, mailID int64) bool {
	this.EnsureUserMail(owner, mailID)
	succ1 := cache.MailCache.DelUserMail(owner, mailID)
	succ2 := cache.MailCache.DelMailCache(mailID)
	if succ1 && succ2 {
		lpc.DBServiceProxy.Delete(&db.DBMail{ID: mailID}, db.DatabaseHandler)
		return true
	}
	return false
}

func (this *mailManager) RemoveMails(owner int32, mailIDs []int64) {
	for _, mailID := range mailIDs {
		if !this.RemoveMail(owner, mailID) {
			continue
		}
	}
}

func (this *mailManager) GetAllMail(owner int32) []*protocol.Mail {
	var results []*protocol.Mail
	mailIDs := cache.MailCache.GetUserMails(owner)
	if mailIDs == nil || len(mailIDs) == 0{
		return results
	}
	for _, mailID := range mailIDs {
		mail :=  cache.MailCache.GetMailCache(int64(mailID))
		if mail != nil {
			results = append(results, mail)
			continue
		}
		cache.MailCache.DelUserMail(owner, int64(mailID))
	}
	return results
}

func (this *mailManager) GetMail(uid int32, limitTime int64, offset int64, count int32) []*protocol.Mail {
	var results []*protocol.Mail

	return results
}