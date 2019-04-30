/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/9/11
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package config

import (
	"flag"
	"io/ioutil"
	"encoding/json"
	"os"
	"aliens/log"
)

var (
	Debug = false
	Tag   = ""
	Root  = ""
	LogDir = ""
)

func init() {
	flag.BoolVar(&Debug, "debug", false, "debug flag")
	flag.StringVar(&Tag, "tag", "", "statistic tag")
	flag.StringVar(&Root, "c", "", "configuration root path")
	flag.StringVar(&LogDir, "log", "", "log root path")
	flag.Parse()
	if Root != "" {
		Root = Root + string(os.PathSeparator)
	}
}

func LoadConfigData(path string, config interface{}) {
	LoadConfigDataEx(path, config,true)
}

func GetConfigPath(path string) string {
	return Root + path
}

func LoadConfigDataEx(path string, config interface{}, fatal bool) {
	if config == nil {
		return
	}
	path = Root + path
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if fatal {
			log.Fatal("config file %v  is not found", path)
		}
		return
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		if fatal {
			log.Fatal("load config %v err %v", path, err)
		}
	}
}


