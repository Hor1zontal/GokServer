/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/5/10
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package conf

import (
	"gok/config"
)

var Config struct {
	Enable			bool
	ES 				ESConfig
}

type ESConfig struct {
	Name string
	Url string
	Host string
	Username string
	Password string
}


func init() {
	config.LoadConfigDataEx("conf/gok/statistics/server.json", &Config, false)
	if Config.ES.Name == "" {
		Config.ES.Name = "gok"
	}
}
