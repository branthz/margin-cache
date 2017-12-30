package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/branthz/utarrow/lib/log"
	"github.com/margin-cache/hashmap"
)

var (
	mlog       *log.Logger
	CacheSet   [KEYHASHMAX]*hashmap.Cache
	GstartTime = time.Now().Unix()
)

func main() {
	var err error
	mlog, err = log.New("", CFV.loglevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	mlog.Info("commonCache start")

	runtime.GOMAXPROCS(runtime.NumCPU())
	for i := 0; i < KEYHASHMAX; i++ {
		CacheSet[i] = hashmap.New(hashmap.NoExpiration, hashmap.DefaultCleanUpInterval)
	}

	//gcacher = hashmap.New(hashmap.NoExpiration, hashmap.DefaultCleanUpInterval)
	//fd, _ := os.Create("./aaa.pprof")
	//pprof.StartCPUProfile(fd)
	//pprof.WriteHeapProfile(fd)
	tListenRequest()
	//time.Sleep(time.Second * 120)
	//pprof.StopCPUProfile()
	//fd.Close()
}
