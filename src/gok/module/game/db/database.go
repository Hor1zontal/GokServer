package db

import (
	"aliens/database"
	"aliens/database/mongo"
	"gok/module/game/conf"
)

var Database database.IDatabase = &mongo.Database{}
var DatabaseHandler = Database.GetHandler()
//var UpdateHandler = database.NewDBUpdateHandler(DatabaseHandler, &database.DBUpdateConfig{30, 2000, 500})

func Init() {
	if conf.Server.Database.Name == "" {
		conf.Server.Database.Name = "gok"
	}
	err := Database.Init(conf.Server.Database)
	if err != nil {
		panic(err)
	}

	DatabaseHandler.EnsureTable("role", &DBRole{})
	DatabaseHandler.EnsureTable("user_message", &DBMessage{})
}

func Close() {
	//Database.Close()()
}


