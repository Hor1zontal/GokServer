/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/11/14
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)


type UserFilter struct {
	UserID   int32     `json:"uid"`              //用户id
	Time     time.Time `json:"time"`
	Content  string    `json:"-"`
}

type LogFilter struct {
	userLogs map[int32]*UserFilter
}


func (this *LogFilter) ReadFile(filePath string, start time.Time, end time.Time) {

	srcFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer srcFile.Close()
	rd := bufio.NewReader(srcFile)

	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		data := &UserFilter{}
		err1 := json.Unmarshal([]byte(line), data)
		if err != nil {
			fmt.Println(err1)
			return
		}
		oldData := this.userLogs[data.UserID]
		
		if data.Time.Before(start) || data.Time.After(end) {
			continue
		}

		if oldData == nil || data.Time.After(oldData.Time) {
			data.Content = line
			this.userLogs[data.UserID] = data
		}

		fmt.Println(line)
	}
}

func (this *LogFilter) WriteFile(filePath string) {
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("create map file error: %v\n", err)
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, v := range this.userLogs {
		fmt.Fprint(w, v.Content)
	}
	w.Flush()
}


func (this *LogFilter) Parser(dirPath string, prefixs []string, distPath string, start time.Time, end time.Time) {
	dir_list, e := ioutil.ReadDir(dirPath)
	if e != nil {
		fmt.Println("read dir error")
		return
	}
	for _, v := range dir_list {
		for _, prefix := range prefixs {
			if !v.IsDir() && strings.HasPrefix(v.Name(), prefix) {
				this.ReadFile(dirPath + "/" + v.Name(), start, end)
			}
		}
	}
	this.WriteFile(distPath)
}

func main() {
	var now = ""
	var day = ""
	var inPath = ""
	var outPath = ""
	var logName = ""

	flag.StringVar(&now, "now", "", "time now")
	flag.StringVar(&day, "day", "", "interval day")
	flag.StringVar(&inPath,"inPath", "", "in dir path")
	flag.StringVar(&outPath,"outPath", "", "out dir path")
	flag.StringVar(&logName, "logName", "", "log name")

	flag.Parse()
	filter := &LogFilter{userLogs:make(map[int32]*UserFilter)}

	//filter.Parser("/Users/hejialin/Documents/aliens/gok/log", "logout.log.2018", "/Users/hejialin/Documents/aliens/gok/logout.log.all")
	//获取前一天的日期
	//t := time.Now()
	//now := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	t, err := time.Parse("20060102", now)
	//当天零点的本地时间
	nowTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	if err != nil {
		fmt.Println(err.Error())
	}
	var prefixs []string
	result, _ := strconv.Atoi(day)
	logNum := result + 1
	for i:=1 ; i <= logNum ; i ++ {
		prefixs = append(prefixs, logName + ".log." + getBeforeDate(nowTime, i))
	}
	fmt.Println(prefixs)
	outDist := outPath + "/" + logName + ".log." + getBeforeDate(nowTime, 1) + "_" +day
	start := nowTime.AddDate(0, 0, -result)
	end := nowTime

	fmt.Println("start time:", start)
	fmt.Println("end time:", end)
	filter.Parser(inPath, prefixs, outDist, start, end)
}

func getBeforeDate(t time.Time, day int) string {
	return  t.AddDate(0, 0, -day).Format("20060102")
}