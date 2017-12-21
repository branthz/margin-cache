package main

import (
	"fmt"

	"github.com/branthz/utarrow/zconfig"
)

type defaultOp struct {
	outPort  int
	loglevel int
}

var Param zconfig.Result
var CFV *defaultOp

func newConfig() *defaultOp {
	op := new(defaultOp)
	op.loglevel = 4
	return op
}

func uploadParameters() error {
	var err error
	CFV = newConfig()

	CFV.loglevel, err = Param.GetInt("logLevel")
	if err != nil {
		return err
	}

	CFV.outPort, err = Param.GetInt("tcpListenPort")
	if err != nil {
		return err
	}

	//get this server inner's addr
	fmt.Printf("tcp listen on port:%d\n", CFV.outPort)

	return nil
}
