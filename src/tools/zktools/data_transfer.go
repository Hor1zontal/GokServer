/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/9/13
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package zktools

import (
	"github.com/name5566/leaf/log"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

//迁移节点数据数据

func Transfer(oldAddress string, newAddress string, oldPath string, newPath string) {
	oldCon, _, err1 := zk.Connect([]string{oldAddress}, time.Second)
	newCon, _, err2 := zk.Connect([]string{newAddress}, time.Second)
	if err1 != nil {
		log.Error("old connection error : %v", err1)
		return
	}

	if err2 != nil {
		log.Error("new connection error : %v", err2)
		return
	}
	copy(oldCon, newCon, oldPath, newPath)
}

func Update(address string, path string, data []byte) {
	con, _, err1 := zk.Connect([]string{address}, time.Second)
	if err1 != nil {
		log.Error("connection error : %v", err1)
		return
	}

	UpdateByPath(con, path, data)
}

func Connect(address string) (*zk.Conn, error) {
	con, _, err := zk.Connect([]string{address}, time.Second)
	if err != nil {
		return nil, err
	}
	return con, nil
}

func Create(con *zk.Conn, path string) {
	result, _, _ := con.Exists(path)
	if !result {
		con.Create(path, []byte(""), 0, zk.WorldACL(zk.PermAll))
	}
}

func UpdateByPath(con *zk.Conn, path string, data []byte) {
	result, _, _ := con.Exists(path)
	if !result {
		_, err2 := con.Create(path, data, 0, zk.WorldACL(zk.PermAll))
		if err2 != nil {
			log.Debug("%v", err2)
			return
		}
	}
	_, err := con.Set(path, data, -1)
	if err != nil {
		log.Error("%v", err)
		return
	}
	log.Release("update success : %v", path)
}

func copy(oldConn *zk.Conn, newConn *zk.Conn, oldPath string, newPath string) bool {
	data, _, err1 := oldConn.Get(oldPath)
	if err1 != nil {
		log.Error("get data error : %v", err1)
		return false
	}
	UpdateByPath(newConn, newPath, data)

	childNames, _, err3 := oldConn.Children(oldPath)
	if err3 != nil {
		log.Error("get children error : %v", err3)
		return false
	}
	for _, name := range childNames {
		copy(oldConn, newConn, oldPath+"/"+name, newPath+"/"+name)
	}
	return true
}
