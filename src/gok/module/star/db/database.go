package db

import (
	"aliens/database"
	"aliens/database/mongo"
	"gok/module/star/conf"
)

var Database database.IDatabase = &mongo.Database{}
var DatabaseHandler = Database.GetHandler()
//var UpdateHandler = database.NewDBUpdateHandler(DatabaseHandler, &database.DBUpdateConfig{60, 2000, 500})

func Init() {
	if conf.Server.Database.Name == "" {
		conf.Server.Database.Name = "gok_star"
	}
	err := Database.Init(conf.Server.Database)
	if (err != nil) {
		panic(err)
	}
	DatabaseHandler.EnsureTable("star", &DBStar{})
}

func Close() {
	//UpdateHandler.Close()
	//Database.Close()()
}

