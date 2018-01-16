package hashmap

import(
	"sync"
	"time"
	"fmt"
)

type Hcache struct{
	mp map[string]*cache
	mu sync.RWMutex 
}

func newHcache()*Hcache{
	st :=new(Hcache)
	st.mp=make(map[string]*cache)
	return st
}

func (h *Hcache)hdel(k,f string) {
	h.mu.RLock()
	c,found:=h.mp[k]
	h.mu.RUnlock()
	if !found{
		return
	}
	c.Delete(f)
	return
}

func (h *Hcache)hdes(k string) {
	h.mu.Lock()
	delete(h.mp,k)
	h.mu.Unlock()
	return
}

func (h *Hcache)get(k string)(*cache,bool) {
	h.mu.RLock()
	v,found:=h.mp[k]
	h.mu.RUnlock()
	return v,found
}

func (h *Hcache)hget(k,f string)(interface{},error) {
	h.mu.RLock()
	c,found:=h.mp[k]
	if !found{
		h.mu.RUnlock()
		return nil, fmt.Errorf("not find the key:%s",k)
	}
	h.mu.RUnlock()
	v,found:=c.Get(f)	
	if !found{
		return nil, fmt.Errorf("not find the field:%s",f)
	}
	return v,nil
}

func (h *Hcache)getOrCreate(k string) (c *cache) {
	var found bool
	h.mu.RLock()
	c,found=h.mp[k]
	if !found{
		c=newCache(DefaultExpiration)
		h.mp[k]=c
	}
	h.mu.RUnlock()
	return
}

func (h *Hcache)hset(k,f string,x interface{},fe time.Duration){			
	c:=h.getOrCreate(k)
	c.Set(f,x,fe)
	return
}

func (h *Hcache)hexist(k,f string) bool{
	h.mu.RLock()
	c,found:=h.mp[k]
	if !found{
		h.mu.RUnlock()
		return false
	}
	_,found=c.Get(f)
	if !found{
		return false
	}
	return true
}

