/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/8/1
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package analysis

import (
	_ "expvar"
	_ "net/http/pprof"
	"gok/constant"
	"net/http"
	"os"
	"runtime/trace"
	"fmt"
	"runtime"
)

func Init() {
	if constant.ANALYSIS {
		//关闭GC
		//debug.SetGCPercent(-1)
		//运行trace
		http.HandleFunc("/start", traces)
		//停止trace
		http.HandleFunc("/stop", traceStop)
		//手动GC
		http.HandleFunc("/gc", gc)
		go http.ListenAndServe(":6060", nil)
	}
}



//手动GC
func gc(w http.ResponseWriter, r *http.Request) {
	runtime.GC()
	w.Write([]byte("StartGC"))
}


//运行trace
func traces(w http.ResponseWriter, r *http.Request){
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}


	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	w.Write([]byte("TrancStart"))
	fmt.Println("StartTrancs")
}

//停止trace
func traceStop(w http.ResponseWriter, r *http.Request){
	trace.Stop()
	w.Write([]byte("TrancStop"))
	fmt.Println("StopTrancs")
}


