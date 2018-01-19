package proxy

import (
	"net"
	"sync"
	"time"

	"github.com/branthz/margin-cache/cmargin"
	"github.com/branthz/margin-cache/common/log"
)

const (
	MonitorInterval = 1e9
)

type remote struct {
	mu         sync.Mutex
	addr       string
	inactive   bool
	needrehash bool
}

func (r *remote) inactivate() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.inactive = true
}

func (r *remote) isActive() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return !r.inactive
}

func (r *remote) tryReactivate() error {
	conn, err := net.Dial("tcp", r.addr)
	if err != nil {
		return err
	}
	conn.Close()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.inactive = false
	return nil
}

func (c *config) rehashing(r *remote) {
	cl := cmargin.Client{
		Addr:        r.addr,
		MaxPoolSize: 1,
	}
	defer cl.Close()
	//TODO add an command to  get all key and value
	keys, err := cl.Keys("*")
	if err != nil {
		log.Errorln(err)
		return
	}

	var end string
	var connpool = make(map[string]cmargin.Client)
	for _, v := range keys {
		end, err = c.consis.Get(v)
		if err != nil {
			log.Error("rehashing:consistent get:%v", err)
			continue
		}
		if end == r.addr {
			continue
		} else {
			//TODO  set key
			var cll cmargin.Client
			var ok bool
			if cll, ok = connpool[end]; !ok {
				cll = cmargin.Client{
					Addr:        r.addr,
					MaxPoolSize: 1,
				}
				connpool[end]=cll
			}
			//TODO reset key
				
		}
	}
}

func (tp *config) runMonitor() {
	for {
		select {
		case <-time.After(MonitorInterval):
			for _, r := range tp.endpoints {
				if !r.isActive() {
					go func() {
						if err := r.tryReactivate(); err != nil {
							log.Warn("failed to activate endpoint [%s] due to %v", r.addr, err)
						} else {
							log.Info("activated %s", r.addr)
							r.inactive = false
						}
					}()
				}
				//TODO check rehashing
			}
		}
	}
}
