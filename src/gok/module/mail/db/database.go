package db

import (
	"aliens/database"
	"aliens/database/mongo"
	"gok/module/mail/conf"
)

var Database database.IDatabase = &mongo.Database{}
var DatabaseHandler = Database.GetHandler()

func Init() {
	if conf.Server.Database.Name == "" {
		conf.Server.Database.Name = "gok_mail"
	}
	err := Database.Init(conf.Server.Database)
	if err != nil {
		panic(err)
	}
	DatabaseHandler.EnsureTable("mail", &DBMail{})
}

func Close() {
	//Database.Close()()
}

