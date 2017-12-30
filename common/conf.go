package common

import (
	"flag"
	"fmt"
	"os"

	"github.com/branthz/utarrow/zconfig"
)

const (
	sectionName = "marginCache"
)

type defaultOp struct {
	outPort  int
	loglevel int
}

var (
	filePath *string = flag.String("f", "/etc/marginCache.toml", "keep the config info")
	CFV      *defaultOp
)

func init() {
	flag.Parse()
	NewConfig()
	err := CFV.SetUp()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func NewConfig() {
	CFV := new(defaultOp)
	CFV.loglevel = 4
}

func (o *defaultOp) SetUp() error {
	pram, err := zconfig.Readfile(*filePath, sectionName)
	if err != nil {
		return err
	}
	o.loglevel, err = pram.GetInt("loglevel")
	if err != nil {
		return err
	}
	o.outPort, err = pram.GetInt("tcplistenport")
	if err != nil {
		return err
	}
	fmt.Printf("tcp listen on port:%d\n", o.outPort)
	return nil
}
