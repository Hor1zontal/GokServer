/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/5/10
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package db

import (
	"aliens/database"
	"aliens/database/mongo"
	"gok/module/log/conf"
)

var Database database.IDatabase = &mongo.Database{}
var DatabaseHandler = Database.GetHandler()
//var UpdateHandler = database.NewDBUpdateHandler(DatabaseHandler, &database.DBUpdateConfig{30, 2000, 500})

func Init() {
	if conf.Config.Database.Name == "" {
		conf.Config.Database.Name = "gok_log"
	}
	err := Database.Init(conf.Config.Database)
	if err != nil {
		panic(err)
	}


	DatabaseHandler.EnsureTable("guide_records", &GuideRecord{})
	DatabaseHandler.EnsureTable("item_records", &ItemRecord{})
	DatabaseHandler.EnsureTable("power_records", &PowerRecord{})
	DatabaseHandler.EnsureTable("faith_records", &FaithRecord{})
	DatabaseHandler.EnsureTable("gaypoint_records", &GayPointRecord{})
	DatabaseHandler.EnsureTable("diamond_records", &DiamondRecord{})
	DatabaseHandler.EnsureTable("order_records", &OrderRecord{})
	DatabaseHandler.EnsureTable("login_records", &LoginRecord{})
	DatabaseHandler.EnsureTable("register_records", &RegisterRecord{})
	DatabaseHandler.EnsureTable("daylogin_records", &DayLoginRecord{})
	DatabaseHandler.EnsureTable("dayregiste_records",&DayRegisterRecord{})
	DatabaseHandler.EnsureTable("daycharge_records",&DayChargeRecord{})
}

func Close() {
	//Database.Close()()
}
