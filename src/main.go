package main

import (
	"flag"
	"fmt"
	"os"
	"package/hashmap"
	"package/log"
	"package/zlconfig"
	"runtime"
	"time"
)

var (
	mlog       *log.Logger
	filePath   *string = flag.String("f", "/etc/fusionCache.toml", "keep the config info")
	CacheSet   [KEYHASHMAX]*hashmap.Cache
	GstartTime = time.Now().Unix()
)

func main() {
	flag.Parse()
	var err error
	Param, err = zlconfig.Readfile(*filePath, "fusionCache")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	err = uploadParameters()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
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
	go showStatus()
	//fd, _ := os.Create("./aaa.pprof")
	//pprof.StartCPUProfile(fd)
	//pprof.WriteHeapProfile(fd)
	tListenRequest()
	//time.Sleep(time.Second * 120)
	//pprof.StopCPUProfile()
	//fd.Close()
}
