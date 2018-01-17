// Copyright 2017 The margin Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/branthz/margin-cache/common"
	"github.com/branthz/margin-cache/common/log"
	"github.com/branthz/margin-cache/handle"
	"github.com/branthz/margin-cache/hashmap"
)

func main() {
	common.Init()
	log.Info("marginCache start...")
	hashmap.DBSetup(hashmap.NoExpiration, hashmap.DefaultCleanUpInterval)
	handle.Start(common.CFV.Outport)
}
