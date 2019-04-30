package db

import (
	"aliens/database"
	"aliens/database/mongo"
	"gok/module/trade/conf"
	"gok/service/msg/protocol"
)

var Database database.IDatabase = &mongo.Database{}
var DatabaseHandler = Database.GetHandler()

func Init() {
	if conf.Server.Database.Name == "" {
		conf.Server.Database.Name = "gok_trade"
	}
	err := Database.Init(conf.Server.Database)
	if err != nil {
		panic(err)
	}
	DatabaseHandler.EnsureTable("sale", &UserSale{})
	DatabaseHandler.EnsureTable("goods", &UserGoods{})
	DatabaseHandler.EnsureTable("item_help", &protocol.ItemHelp{})
	DatabaseHandler.EnsureTable("item_help_history", &ItemHelpHistory{})
}

func Close() {
	//UpdateHandler.Close()
	//Database.Close()()
}
