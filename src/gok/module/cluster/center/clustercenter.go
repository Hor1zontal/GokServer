/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/8/16
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package center

import (
	"aliens/center/core"
	"gok/module/cluster/conf"
	"gok/module/cluster/hashring"
)

var ClusterCenter *core.ServiceCenter = &core.ServiceCenter{} //服务中心

var ConfigCenter *core.ConfigCenter = &core.ConfigCenter{} //配置中心

var StarHashring = hashring.NewServiceListener("star")

var TradeHashring = hashring.NewServiceListener("trade")




func Init() {
	ClusterCenter.SetLBS(conf.Config.LBS)
	ClusterCenter.ConnectCluster(conf.Config.ZKServers, 10, conf.Config.ZKName)
	ClusterCenter.AddServiceListener(StarHashring)
	ClusterCenter.AddServiceListener(TradeHashring)
	ConfigCenter.StartCluster(conf.Config.ZKServers, 10, conf.Config.ZKName)
}

func Close() {
	if ClusterCenter != nil {
		ClusterCenter.Close()
	}
	if ConfigCenter != nil {
		ConfigCenter.Close()
	}
}

func IsMaster(serviceType string) bool {
	service := ClusterCenter.GetMasterService(serviceType)
	if service == nil {
		return false
	}
	return service.GetID() == GetServerNode()
}

func GetServerNode() string {
	return conf.Config.Node
}
