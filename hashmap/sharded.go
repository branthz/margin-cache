// Copyright 2017 The margin Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package hashmap provides an key/value store.
package hashmap

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	insecurerand "math/rand"
	"os"
	"time"
)

// Dbs :This is multiple cashes, namely by
// preventing write locks of the entire cache when an item is added. As of the
// time of writing, the overhead of selecting buckets results in cache
// operations being about twice as slow as for the standard cache with small
// total cache sizes, and faster for larger ones.
//
// See cache_test.go for a few benchmarks.
type Dbs struct {
	*shardedCache
}

const (
	//DefaultCleanUpInterval clean the cache expied items
	DefaultCleanUpInterval time.Duration = 60 * 1e9
)

var (
	defaultDbs *Dbs
)

type shardedCache struct {
	seed    uint32
	m       uint32
	cs      []*cache
	h       *Hcache
	janitor *shardedJanitor
}

// djb2 with better shuffling. 5x faster than FNV with the hash.Hash overhead.
func djb33(seed uint32, k string) uint32 {
	var (
		l = uint32(len(k))
		d = 5381 + seed + l
		i = uint32(0)
	)
	// Why is all this 5x faster than a for loop?
	if l >= 4 {
		for i < l-4 {
			d = (d * 33) ^ uint32(k[i])
			d = (d * 33) ^ uint32(k[i+1])
			d = (d * 33) ^ uint32(k[i+2])
			d = (d * 33) ^ uint32(k[i+3])
			i += 4
		}
	}
	switch l - i {
	case 1:
	case 2:
		d = (d * 33) ^ uint32(k[i])
	case 3:
		d = (d * 33) ^ uint32(k[i])
		d = (d * 33) ^ uint32(k[i+1])
	case 4:
		d = (d * 33) ^ uint32(k[i])
		d = (d * 33) ^ uint32(k[i+1])
		d = (d * 33) ^ uint32(k[i+2])
	}
	return d ^ (d >> 16)
}

func (sc *shardedCache) bucket(k string) *cache {
	return sc.cs[djb33(sc.seed, k)%sc.m]
}

func (sc *shardedCache) Set(k string, x interface{}, d time.Duration) {
	sc.bucket(k).Set(k, x, d)
}

func (sc *shardedCache) Add(k string, x interface{}, d time.Duration) error {
	return sc.bucket(k).Add(k, x, d)
}

func (sc *shardedCache) Replace(k string, x interface{}, d time.Duration) error {
	return sc.bucket(k).Replace(k, x, d)
}

func (sc *shardedCache) Get(k string) (interface{}, bool) {
	return sc.bucket(k).Get(k)
}

func (sc *shardedCache) Getallkey(buff *bytes.Buffer) (int, error) {
	var cn, n int
	var err error
	for _, v := range sc.cs {
		n, err = v.Getallkey(buff)
		if err != nil {
			return cn, err
		}
		cn += n
	}
	return cn, err
}

func (sc *shardedCache) Increment(k string, n int64) error {
	return sc.bucket(k).Increment(k, n)
}

func (sc *shardedCache) IncrementInt64(k string, n int64) (int64, error) {
	return sc.bucket(k).IncrementInt64(k, n)
}

func (sc *shardedCache) IncrementFloat(k string, n float64) error {
	return sc.bucket(k).IncrementFloat(k, n)
}

func (sc *shardedCache) Decrement(k string, n int64) error {
	return sc.bucket(k).Decrement(k, n)
}

func (sc *shardedCache) DecrementInt64(k string, n int64) (int64, error) {
	return sc.bucket(k).DecrementInt64(k, n)
}
func (sc *shardedCache) Delete(k string) {
	sc.bucket(k).Delete(k)
}

func (sc *shardedCache) DeleteExpired() {
	for _, v := range sc.cs {
		v.DeleteExpired()
	}
}

func (sc *shardedCache) Hget(k, f string) (interface{}, error) {
	return sc.h.hget(k, f)
}

func (sc *shardedCache) Hset(k, f string, x interface{}, d time.Duration) {
	sc.h.hset(k, f, x, d)
	return
}

func (sc *shardedCache) Hexist(k, f string) bool {
	return sc.h.hexist(k, f)
}

func (sc *shardedCache) Hdel(k, f string) {
	sc.h.hdel(k, f)
	return
}

func (sc *shardedCache) Hmset(k string, pairs [][]byte) {
	for i := 0; i < len(pairs); i = i + 2 {
		sc.h.hset(k, string(pairs[i]), pairs[i+1], NoExpiration)
	}
}

func (sc *shardedCache) Hdestroy(k string) {
	sc.h.hdes(k)
	return
}

func (sc *shardedCache) Hmget(k string, pairs [][]byte) (data [][]byte, err error) {
	c, found := sc.h.get(k)
	if !found {
		err = fmt.Errorf("no find the key:%s", k)
		return
	}

	for i := 0; i < len(pairs); i++ {
		v, found := c.Get(string(pairs[i]))
		if !found {
			err = fmt.Errorf("no find the field:%s", string(pairs[i]))
			return
		}
		data = append(data, v.([]byte))
	}
	return
}

func (sc *shardedCache) Hgetall(k string, buf *bytes.Buffer) error {
	c, ok := sc.h.get(k)
	if !ok {
		return fmt.Errorf("no find key:%s", k)
	}
	err := c.Getall(buf)
	return err
}

// Returns the items in the cache. This may include items that have expired,
// but have not yet been cleaned up. If this is significant, the Expiration
// fields of the items should be checked. Note that explicit synchronization
// is needed to use a cache and its corresponding Items() return values at
// the same time, as the maps are shared.
func (sc *shardedCache) Items() []map[string]Item {
	res := make([]map[string]Item, len(sc.cs))
	for i, v := range sc.cs {
		res[i] = v.Items()
	}
	return res
}

func (sc *shardedCache) Flush() {
	for _, v := range sc.cs {
		v.Flush()
	}
}

type shardedJanitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *shardedJanitor) Run(sc *shardedCache) {
	j.stop = make(chan bool)
	tick := time.Tick(j.Interval)
	for {
		select {
		case <-tick:
			sc.DeleteExpired()
		case <-j.stop:
			return
		}
	}
}

func stopShardedJanitor(sc *Dbs) {
	sc.janitor.stop <- true
}

func runShardedJanitor(sc *shardedCache, ci time.Duration) {
	j := &shardedJanitor{
		Interval: ci,
	}
	sc.janitor = j
	go j.Run(sc)
}

func newShardedCache(n int, de time.Duration) *shardedCache {
	max := big.NewInt(0).SetUint64(uint64(math.MaxUint32))
	rnd, err := rand.Int(rand.Reader, max)
	var seed uint32
	if err != nil {
		os.Stderr.Write([]byte("WARNING: newShardedCache failed to read from the system CSPRNG (/dev/urandom or equivalent.) .Continuing with an insecure seed.\n"))
		seed = insecurerand.Uint32()
	} else {
		seed = uint32(rnd.Uint64())
	}
	hc := newHcache()
	sc := &shardedCache{
		seed: seed,
		m:    uint32(n),
		cs:   make([]*cache, n),
		h:    hc,
	}
	for i := 0; i < n; i++ {
		c := &cache{
			defaultExpiration: de,
			items:             map[string]Item{},
			id:                uint32(i),
		}
		sc.cs[i] = c
	}
	return sc
}

//Maxbuckets define buckets count.
const Maxbuckets = 30

//DBSetup init dbs
func DBSetup(defaultExpiration, cleanupInterval time.Duration) *Dbs {
	if defaultDbs != nil {
		return defaultDbs
	}
	if defaultExpiration == 0 {
		defaultExpiration = -1
	}
	sc := newShardedCache(Maxbuckets, defaultExpiration)
	defaultDbs = &Dbs{sc}
	//if cleanupInterval > 0 {
	//	runShardedJanitor(sc, cleanupInterval)
	//	runtime.SetFinalizer(defaultDbs, stopShardedJanitor)
	//}
	return defaultDbs
}

// GetDB returns the defaultDBs
func GetDB() *Dbs {
	if defaultDbs != nil {
		return defaultDbs
	}
	return DBSetup(NoExpiration, DefaultCleanUpInterval)
}
