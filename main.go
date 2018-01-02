package main

import (
	"github.com/branthz/margin-cache/common"
	"github.com/branthz/margin-cache/common/log"
	"github.com/branthz/margin-cache/handle"
	"github.com/branthz/margin-cache/hashmap"
)

func main() {
	common.Init()
	log.Info("commonCache start")
	dbs := hashmap.DBSetup(hashmap.NoExpiration, hashmap.DefaultCleanUpInterval)

	//for i := 0; i < KEYHASHMAX; i++ {
	//	CacheSet[i] = hashmap.New(hashmap.NoExpiration, hashmap.DefaultCleanUpInterval)
	//}

	//gcacher = hashmap.New(hashmap.NoExpiration, hashmap.DefaultCleanUpInterval)
	//fd, _ := os.Create("./aaa.pprof")
	//pprof.StartCPUProfile(fd)
	//pprof.WriteHeapProfile(fd)
	handle.TListen(dbs)
	//time.Sleep(time.Second * 120)
	//pprof.StopCPUProfile()
	//fd.Close()
}
