/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/9/5
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package app

import (
	"github.com/name5566/leaf/module"
	"aliens/log"
	"github.com/name5566/leaf"
	"math/rand"
	"time"
	"gok/module/statistics/elastics"
	"gok/config"
)



func Run(mods ...module.Module) {
	rand.Seed(time.Now().UnixNano())
	log.Init(config.Debug, config.Tag, config.LogDir)
	log.Info("==============================version:1.0.0.2018102901==============================")
	elastics.Init(config.Tag)
	leaf.Run(mods...)
}

