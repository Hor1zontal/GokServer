/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/9/13
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package main

import (
	"flag"
	"os"
	"strings"

	"tools/zktools"
	"github.com/name5566/leaf/log"
	"path/filepath"
	"fmt"
	"io"
	"bufio"
	"github.com/samuel/go-zookeeper/zk"
	"io/ioutil"
	"crypto/md5"
	"bytes"
	"encoding/json"
)

//var (
//	from 		string
//	to 			string
//	fromPath 	string
//	toPath 		string
//)

var (
	address         string
	jsPath          string
	jsonPath        string
	node            string
	config			string

	checkFilePath   string

	updateNames     []string
	isChange		bool
	checkMd5Mapping map[string][]byte
	unUpdateFiles	[]string
)

func main() {
	//flag.StringVar(&from, "from", "","zookeeper copy from address")
	//flag.StringVar(&to,"to", "", "zookeeper copy to address")
	//flag.StringVar(&fromPath,"fromPath", "", "fromPath")
	//flag.StringVar(&toPath,"toPath", "", "node name")
	//flag.Parse()
	//
	//Transfer(from, to, fromPath, toPath)

	//zookeeper.Transfer("120.77.213.111:2182", "39.108.220.37:2181", "/gok")

	//result, _ := ip.GetLocalPublicIpUseDnspod()
	//log.Debug(result)

	flag.StringVar(&address,"address", "", "zookeeper address")
	flag.StringVar(&jsPath,"jsPath", "", "zookeeper file")
	flag.StringVar(&jsonPath, "jsonPath", "", "zookeeper file")
	flag.StringVar(&node,"node", "", "node name")
	flag.StringVar(&checkFilePath,"checkFilePath", "./check.json", "md5Path")
	flag.StringVar(&config, "blackList", "","")

	flag.Parse()
	uploadWithConfig(address, jsPath, node, jsonPath, checkFilePath, config)
}

func uploadWithConfig(address string, jsPath string, node string, jsonPath string, checkFilePath string, unUpdateFile string) {
	// 初始化部分
	zkConn, err := zktools.Connect(address)
	if err != nil {
		log.Debug("zk Conn error: %v", err.Error())
		return
	}
	nodePath := node
	strArr := strings.Split(nodePath, "/")
	strArr = append(strArr[:0], strArr[1:]...)
	tempPath := ""
	for i := 0; i < len(strArr); i++ {
		tempPath += "/" + strArr[i]
		zktools.Create(zkConn, tempPath)
	}

	f, err := os.Open(checkFilePath)
	if err != nil {
		log.Debug("os open checkFile error: %v", err.Error())
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Debug("data read error: %v", err.Error())
	}
	checkMd5Mapping = make(map[string][]byte)
	err = json.Unmarshal(data, &checkMd5Mapping)
	if err != nil{
		println("%v",err.Error())
	}

	loadUnUpdateFiles(unUpdateFile)

	// js格式上传
	if jsPath != "" {
		files, err := filepath.Glob(jsPath )
		if err != nil {
			log.Debug("jsPath open error: %v", err.Error())
		}
		if files == nil || len(files) == 0{
			log.Debug("there is no js file in this directory: %v",jsPath)
		} else {
			log.Debug("update js file...")
			for i := 0; i < len(files); i++ {
				log.Debug("update file path: %v", files[i])
				updateJsFile(files[i], nodePath, zkConn)
			}
		}
	}
	// json格式上传
	if jsonPath != "" {
		files, err := ioutil.ReadDir(jsonPath)
		//files, err := filepath.Glob(jsonPath )

		if err != nil {
			log.Debug("jsonPath open error: %v", err)
		}
		if files == nil || len(files) == 0{
			log.Debug("there is no json file in this directory: %v", jsonPath)
		} else {
			log.Debug("update json file ... ")
			for i := 0; i < len(files); i++ {
				log.Debug("update file path: %v", jsonPath + files[i].Name())
				updateJsonFile(jsonPath + files[i].Name(), nodePath, zkConn)
			}
		}
	}

	if isChange {
		data, err1 := json.Marshal(&checkMd5Mapping)
		if err1 != nil {
			log.Debug("map marshal error: %v", err1)
		}
		err2 := ioutil.WriteFile(checkFilePath, data, 0644)
		if err2 != nil {
			log.Error("%v", err2.Error())
		}
	}

}

func loadUnUpdateFiles(unUpdateFile string) {
	f, err1 := os.Open(unUpdateFile)
	if err1 != nil {
		log.Debug("os open unUpdateFile error: %v", err1.Error())
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if line == "\n"  {
			//log.Debug("空行")
			continue
		}
		line = strings.Replace(line, "\r", "", -1)
		line = strings.Replace(line, "\n", "", -1)
		unUpdateFiles = append(unUpdateFiles, line)
	}
}

func isAbleUpdate(tableName string) bool {
	for _, name := range unUpdateFiles {
		if tableName == name {
			log.Release("%v has been ignore",name)
			return false
		}
	}
	return true
}

func updateJsFile(filePath string, nodePath string, zkConn *zk.Conn) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Debug("os open error: %v", err.Error())
		return
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}

		var strLine, name, value string
		strLine = string(line)
		if strLine == "\n" || strLine == "\r\n" {
			//log.Debug("空行")
			continue
		}
		strArr := strings.Split(strLine, " ")
		if len(strArr) >= 2 {
			name = strArr[0]                     //第一个字段是表名
			value = strings.Join(strArr[1:], "") //第二个以上的字段有可能被拆分多个,就需要连接到一起
		} else {
			fmt.Println("##error##", strArr[:1], "表格解析失败，请检查")
		}
		if name == "" {
			log.Debug("table name is nil")
			continue
		}
		name = strings.ToLower(name)
		path := nodePath + "/" + name
		if !isAbleUpdate(name) {
			continue
		}
		if !ensureNameExist(name) {
			updateNames = append(updateNames, name)
		}
		updateData := []byte(value)
		md5Data := convertToMd5(updateData)
		if !checkMd5(name, md5Data) {
			zktools.UpdateByPath(zkConn, path, updateData)

			updateMd5(name, md5Data)
			isChange = true
		}
	}
}

func convertToMd5(updateData []byte) []byte {
	h :=md5.New()
	h.Write(updateData) // 需要加密的字符串
	md5Data := h.Sum(nil)
	return md5Data
}

func checkMd5(checkName string, md5Data []byte) bool {
	if bytes.Equal(checkMd5Mapping[checkName], md5Data) {
		log.Release("There is no change in %v",checkName)
		return true
	}
	return false
}

func updateMd5(name string, md5Data []byte) {
	//fmt.Printf("%s\n", hex.EncodeToString(cipherStr)) // 输出加密结果
	checkMd5Mapping[name] = md5Data
}

func updateJsonFile(fileName string, nodePath string, zkConn *zk.Conn) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Error("os open error: %v", err.Error())
		return
	}
	defer f.Close()

	jsonData, err := ioutil.ReadAll(f)
	if err != nil {
		log.Error("json read error: %v", err.Error())
	}
	jsonName := ""
	for i := 0; i < 1000000; i++ {
		jsonName = filepath.Base(f.Name())
		strArr := strings.Split(jsonName,".")
		jsonName = strArr[0]
	}
	jsonName = strings.ToLower(jsonName)
	path := nodePath + "/" + jsonName
	if !isAbleUpdate(jsonName) {
		return
	}

	updateData := jsonData
	if !json.Valid(updateData) {
		log.Error("%v is not json file", jsonName)
		//fmt.Println("##ERROR##", jsonName, "is not a json file")
		return
	}
	md5Data := convertToMd5(updateData)
	if !ensureNameExist(jsonName) {
		updateNames = append(updateNames, jsonName)
	}
	if !checkMd5(jsonName, md5Data) {
		zktools.UpdateByPath(zkConn, path, updateData)
		updateMd5(jsonName, md5Data)
		isChange = true
	}
}

func ensureNameExist(name string) bool {
	for _, updateName := range updateNames {
		if name == updateName {
			fmt.Println("##WARNING##", name, "is already update")
			return true
		}
	}
	return false
}