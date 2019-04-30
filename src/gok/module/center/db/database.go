package db

import (
	"aliens/database/mongo"
	"aliens/database"
	"gok/module/center/conf"
)

var Database database.IDatabase = &mongo.Database{}
var DatabaseHandler = Database.GetHandler()
//var UpdateHandler = database.NewDBUpdateHandler(DatabaseHandler, &database.DBUpdateConfig{30, 2000, 500})

func Init() {
	if conf.Server.Database.Name == "" {
		conf.Server.Database.Name = "gok_center"
	}
	err := Database.Init(conf.Server.Database)
	if (err != nil) {
		panic(err)
	}

	DatabaseHandler.EnsureTable("notices", &DBNotice{})
}

func Close() {
	//Database.Close()()
}
