// Copyright 2017 The margin Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package common provides app's configure and tools.
package common

import (
	"flag"
	"fmt"
	"os"

	"github.com/branthz/margin-cache/common/log"
	"github.com/branthz/utarrow/zconfig"
)

const (
	sectionName = "marginCache"
)

//AppOp saves the app configure
type AppOp struct {
	Outport  int
	loglevel int
}

var (
	filePath *string = flag.String("f", "/etc/marginCache.toml", "keep the config info")
	CFV      *AppOp
)

//Init the app conf and env
func Init() {
	flag.Parse()
	newConfig()
	err := CFV.SetUp()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	err = log.Setup("", CFV.loglevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func newConfig() {
	CFV = new(AppOp)
	CFV.loglevel = 4
}

// SetUp init config
func (o *AppOp) SetUp() error {
	pram, err := zconfig.Readfile(*filePath, sectionName)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", pram)
	o.loglevel, err = pram.GetInt("logLevel")
	if err != nil {
		return err
	}
	o.Outport, err = pram.GetInt("tcpListenPort")
	if err != nil {
		return err
	}
	fmt.Printf("tcp listen on port:%d\n", o.Outport)
	return nil
}
