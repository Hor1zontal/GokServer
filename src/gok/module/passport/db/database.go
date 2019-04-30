package db

import (
	"gok/module/passport/conf"
	"aliens/database/mongo"
	"aliens/database"
)

var Database database.IDatabase = &mongo.Database{}
var DatabaseHandler = Database.GetHandler()

func Init() {
	if conf.Server.Database.Name == "" {
		conf.Server.Database.Name = "gok_passport"
	}
	err := Database.Init(conf.Server.Database)
	if err != nil {
		panic(err)
	}
	DatabaseHandler.EnsureTable("user", &DBUser{})

}

func Close() {
	//Database.Close()()
}
