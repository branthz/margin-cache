package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type statesVal struct {
	v int64
}

func increment(sv *statesVal) {
	atomic.AddInt64(&sv.v, 1)
}
func decrement(sv *statesVal) {
	atomic.AddInt64(&sv.v, -1)
}
func assign(sv *statesVal, count int64) {
	sv.v = count
}
func addStateVal(sv *statesVal) {
	sv.v++
}

var statsMap = map[string]*statesVal{
	"Fusion.clientSum": &statesVal{0},
	"Fusion.clientNow": &statesVal{0},
	"Fusion.keysNum":   &statesVal{0},
}

func printStats() {
	fmt.Printf("%s:", time.Now().Format(time.RFC3339))
	for k, v := range statsMap {
		fmt.Printf("%s:%d;\t\t", k, v.v)
	}
	fmt.Printf("\n")
}

func showStatus() {
	t3 := time.NewTicker(time.Second * 120) //print statics
	for {
		select {
		case <-t3.C:
			var count int = 0
			for i := 0; i < KEYHASHMAX; i++ {
				count = count + CacheSet[i].ItemCount()
			}
			assign(statsMap["Fusion.keysNum"], int64(count))
			printStats()
		}
	}
}
