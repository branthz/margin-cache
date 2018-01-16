package main

import (
	"github.com/branthz/margin-cache/proxy"
	"github.com/branthz/utarrow/lib/log"
	"github.com/branthz/utarrow/zconfig"
	"flag"
)

var (
	filePath *string = flag.String("f", "/etc/marginCache.toml", "keep the config info")
	backends []string
	selfAddr string 
	mlog *log.Logger 
)

func setup(){
	pram, err := zconfig.Readfile(*filePath, "proxy")
	if err != nil {
		panic(err)
	}
	backends,_=pram.GetArrayStr("backends")
	selfAddr,_=pram.GetString("listen")
	mlog,_=log.New("",log.Debug)
}

func main() {
	setup()
	p := &proxy.Proxy{LocalHost:selfAddr}
	for _,host:=range backends{
		p.AddRoute(host)
	}
	if err := p.Run(); err != nil {
		mlog.Errorln(err)
		return
	}
}
