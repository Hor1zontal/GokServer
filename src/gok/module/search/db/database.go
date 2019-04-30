package db

import (
	"aliens/database"
	"aliens/database/mongo"

	"gok/module/search/conf"
)

var Database database.IDatabase = &mongo.Database{}
var DatabaseHandler = Database.GetHandler()

func Init() {
	if conf.Server.Database.Name == "" {
		conf.Server.Database.Name = "gok_search"
	}
	err := Database.Init(conf.Server.Database)
	if (err != nil) {
		panic(err)
	}
	DatabaseHandler.EnsureTable("star_building_info", &StarBuildingInfo{})
	DatabaseHandler.EnsureTable("star_believer_info", &StarBelieverInfo{})
	DatabaseHandler.EnsureTable("user_index_info", &UserIndex{})


}

func Close() {
	//Database.Close()()
}

