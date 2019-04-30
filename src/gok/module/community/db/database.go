package db

import (
	"aliens/database"
	"aliens/database/mongo"
	"gok/module/community/conf"
)

var Database database.IDatabase = &mongo.Database{}
var DatabaseHandler = Database.GetHandler()

func Init() {
	if conf.Server.Database.Name == "" {
		conf.Server.Database.Name = "gok_community"
	}
	err := Database.Init(conf.Server.Database)
	if err != nil {
		panic(err)
	}

	DatabaseHandler.EnsureTable("moments", &DMoments{})
	DatabaseHandler.EnsureTable("follow", &DFollow{})
}

func Close() {
	//UpdateHandler.Close()
	//Database.Close()()
}
