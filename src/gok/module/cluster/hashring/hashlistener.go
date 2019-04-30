/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/7/13
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package hashring

import "aliens/log"


func NewServiceListener(serviceType string) *ServiceListener {
	return &ServiceListener{
		hashRing:NewHashRing(400),
		serviceType:serviceType,
	}
}

type ServiceListener struct {
	hashRing *HashRing
	serviceType string
}

func (this *ServiceListener) GetServiceType() string {
	return this.serviceType
}

func (this *ServiceListener) GetNode(id string) string {
	return this.hashRing.GetNode(id)
}

func (this *ServiceListener) AddNode(id string) {
	this.hashRing.AddNode(id, 1)
	log.Debug("add %v node %v",this.serviceType, id)
}

func (this *ServiceListener) RemoveNode(id string) {
	this.hashRing.RemoveNode(id)
	log.Debug("remove %v node %v",this.serviceType, id)
}
