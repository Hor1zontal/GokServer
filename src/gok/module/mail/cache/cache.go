package cache

import (
	"gok/module/mail/db"
	basecache "gok/cache"
	"aliens/log"
	"gok/module/mail/conf"
	"time"
	"gok/service/lpc"
)

var MailCache = basecache.NewMailCacheManager()

func Init() {
	MailCache.Init1(conf.Server.RedisAddress, conf.Server.RedisPassword,
		conf.Server.RedisMaxActive, conf.Server.RedisMaxIdle, conf.Server.RedisIdleTimeout)

	if MailCache.SetNX(basecache.FLAG_LOAD_MAIL, 1) {
		log.Debug("start load mail data to redis cache...")
		count := 0
		var mails []*db.DBMail
		err := db.DatabaseHandler.QueryAllLimit(&db.DBMail{}, &mails, 10000, func(data interface{}) bool {
			for _, mail := range mails {
				//过期
				var seconds = int(time.Now().Sub(mail.CreateTime).Seconds())
				if seconds > conf.Server.RedisExpireTime {
					lpc.DBServiceProxy.Delete(&db.DBMail{ID:mail.ID},db.DatabaseHandler)
					continue
				}
				//未过期
				MailCache.SetExpireTime(mail.ID, conf.Server.RedisExpireTime - seconds)
				MailCache.SetUserMail(mail.Owner, mail.ID)
			}
			currLen := len(mails)
			count += currLen
			return currLen == 0
		})
		if err != nil {
			log.Debug("load mail err: %v", err)
		}
		log.Debug("end load mail data to redis cache count:%v", count)
	}

	//if MailCache.SetNX(basecache.FLAG_LOAD_MAIL, 1) {
	//	log.Debug("start load mail data to redis cache...")
	//	var mails []*db.DBMail
	//	db.DatabaseHandler.QueryAll(&db.DBMail{}, &mails)
	//	for _, mail := range mails {
	//		//过期
	//		var seconds = int(time.Now().Sub(mail.CreateTime).Seconds())
	//		if seconds > conf.Server.RedisExpireTime {
	//			lpc.DBServiceProxy.Delete(&db.DBMail{ID:mail.ID},db.DatabaseHandler)
	//			continue
	//		}
	//		//未过期
	//		MailCache.SetExpireTime(mail.ID, conf.Server.RedisExpireTime - seconds)
	//		MailCache.SetUserMail(mail.Owner, mail.ID)
	//	}
	//	log.Debug("end load mail data to redis cache")
	//}
}

func Close() {
	MailCache.Close()
}

