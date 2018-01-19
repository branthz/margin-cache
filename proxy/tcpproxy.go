// Copyright 2017 The margin Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package proxy provides an tcp proxy applied with consistent routers to backends
package proxy

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

// Proxy is a proxy.
type Proxy struct {
	configs   *config
	LocalHost string
	lns       net.Listener
	donec     chan struct{} // closed before err
	err       error         // any error from listening

	// ListenFunc optionally specifies an alternate listen function.
	ListenFunc func(net, laddr string) (net.Listener, error)
}

// Matcher reports whether hostname matches the Matcher's criteria.
type Matcher func(ctx context.Context, hostname string) bool

// equals is a trivial Matcher that implements string equality.
func equals(want string) Matcher {
	return func(_ context.Context, got string) bool {
		return want == got
	}
}

// config contains the proxying state for one listener.
type config struct {
	routes []route
	consis *Consistent
	endpoints []*remote
}

func (c *config) match(b *bufio.Reader) Target {
	key := "hello"
	t, _ := c.consis.GetTarget(key)
	return t
}

func (c *config) addRoute(dst string, tar Target) {
	r:=new(remote)
	r.addr=dst 
	r.inactive=false 
	c.endpoints=append(c.endpoints,r)
	c.consis.Add(dst, tar)
}

func (c *config) removeEnd(id string){
	c.consis.Remove(id)
	for _,v:=range c.endpoints{
		if v.addr !=id {
			v.needrehash = true	
		}	
	}
}

// A route matches a connection to a target.
type route interface {
	match(*bufio.Reader) Target
}

func (p *Proxy) netListen() func(net, laddr string) (net.Listener, error) {
	if p.ListenFunc != nil {
		return p.ListenFunc
	}
	return net.Listen
}

func (p *Proxy) configFor() *config {
	if p.configs == nil {
		p.configs = &config{consis: New()}
	}
	return p.configs
}

func (p *Proxy) addRoute(r route) {
	cfg := p.configFor()
	cfg.routes = append(cfg.routes, r)
}

// AddRoute is generally used as either the only rule (for simple TCP proxies based on consistent hash)
//
func (p *Proxy) AddRoute(dest string) {
	p.configFor().addRoute(dest, To(dest))
	//p.addRoute(fixedTarget{To(dest)})
}

type fixedTarget struct {
	t Target
}

func (m fixedTarget) match(*bufio.Reader) Target {
	return m.t
}

// Run is calls Start, and then Wait.
func (p *Proxy) Run() error {
	if err := p.Start(); err != nil {
		return err
	}
	return p.Wait()
}

// Wait waits for the Proxy to finish running. 
func (p *Proxy) Wait() error {
	<-p.donec
	return p.err
}

// Close closes all the proxy's self-opened listeners.
func (p *Proxy) Close() error {
	p.lns.Close()
	return nil
}

// Start creates a TCP listener
// and starts the proxy.
func (p *Proxy) Start() error {
	if p.donec != nil {
		return errors.New("already started")
	}
	p.donec = make(chan struct{})
	errc := make(chan error, 1)
	ln, err := p.netListen()("tcp", p.LocalHost)
	if err != nil {
		p.Close()
		return err
	}
	p.lns = ln
	go p.serveListener(errc, ln, p.configs)
	go p.awaitFirstError(errc)
	return nil
}

func (p *Proxy) awaitFirstError(errc <-chan error) {
	p.err = <-errc
	close(p.donec)
}

func (p *Proxy) serveListener(ret chan<- error, ln net.Listener, route *config) {
	for {
		c, err := ln.Accept()
		if err != nil {
			ret <- err
			return
		}
		go p.serveConn(c, route)
	}
}

// serveConn runs in its own goroutine and matches c against routes.
// It returns whether it matched purely for testing.
func (p *Proxy) serveConn(c net.Conn, route *config) bool {
	br := bufio.NewReader(c)

	if target := route.match(br); target != nil {
		if n := br.Buffered(); n > 0 {
			peeked, _ := br.Peek(br.Buffered())
			c = &Conn{
				Peeked: peeked,
				Conn:   c,
			}
		}
		target.HandleConn(c)
		return true
	}
	// TODO: hook for this?
	log.Printf("tcpproxy: no routes matched conn %v/%v; closing", c.RemoteAddr().String(), c.LocalAddr().String())
	c.Close()
	return false
}

type Conn struct {
	// Peeked are the bytes that have been read from Conn for the
	// purposes of route matching, but have not yet been consumed
	// by Read calls. It set to nil by Read when fully consumed.
	Peeked []byte

	// Conn is the underlying connection.
	net.Conn
}

func (c *Conn) Read(p []byte) (n int, err error) {
	if len(c.Peeked) > 0 {
		n = copy(p, c.Peeked)
		c.Peeked = c.Peeked[n:]
		if len(c.Peeked) == 0 {
			c.Peeked = nil
		}
		return n, nil
	}
	return c.Conn.Read(p)
}

// Target is what an incoming matched connection is sent to.
type Target interface {
	// HandleConn is called when an incoming connection is
	// matched. 
	HandleConn(net.Conn)
}

func To(addr string) *DialProxy {
	return &DialProxy{Addr: addr}
}

// DialProxy implements Target by dialing a new connection to Addr
// and then proxying data back and forth.
type DialProxy struct {
	// Addr is the TCP address to proxy to.
	Addr string

	KeepAlivePeriod time.Duration

	// If zero, a default is used.
	// If negative, the timeout is disabled.
	DialTimeout time.Duration

	// DialContext optionally specifies an alternate dial function
	// for TCP targets. If nil, the standard
	// net.Dialer.DialContext method is used.
	DialContext func(ctx context.Context, network, address string) (net.Conn, error)

	// OnDialError optionally specifies an alternate way to handle errors dialing Addr.
	OnDialError func(src net.Conn, dstDialErr error)

	ProxyProtocolVersion int
}

// UnderlyingConn returns c.Conn if c of type *Conn,
// otherwise it returns c.
func UnderlyingConn(c net.Conn) net.Conn {
	if wrap, ok := c.(*Conn); ok {
		return wrap.Conn
	}
	return c
}

// HandleConn implements the Target interface.
func (dp *DialProxy) HandleConn(src net.Conn) {
	ctx := context.Background()
	var cancel context.CancelFunc
	if dp.DialTimeout >= 0 {
		ctx, cancel = context.WithTimeout(ctx, dp.dialTimeout())
	}
	dst, err := dp.dialContext()(ctx, "tcp", dp.Addr)
	if cancel != nil {
		cancel()
	}
	if err != nil {
		dp.onDialError()(src, err)
		return
	}
	defer dst.Close()

	if err = dp.sendProxyHeader(dst, src); err != nil {
		dp.onDialError()(src, err)
		return
	}
	defer src.Close()

	if ka := dp.keepAlivePeriod(); ka > 0 {
		if c, ok := UnderlyingConn(src).(*net.TCPConn); ok {
			c.SetKeepAlive(true)
			c.SetKeepAlivePeriod(ka)
		}
		if c, ok := dst.(*net.TCPConn); ok {
			c.SetKeepAlive(true)
			c.SetKeepAlivePeriod(ka)
		}
	}

	errc := make(chan error, 1)
	go proxyCopy(errc, src, dst)
	go proxyCopy(errc, dst, src)
	<-errc
}

func (dp *DialProxy) sendProxyHeader(w io.Writer, src net.Conn) error {
	switch dp.ProxyProtocolVersion {
	case 0:
		return nil
	case 1:
		var srcAddr, dstAddr *net.TCPAddr
		if a, ok := src.RemoteAddr().(*net.TCPAddr); ok {
			srcAddr = a
		}
		if a, ok := src.LocalAddr().(*net.TCPAddr); ok {
			dstAddr = a
		}

		if srcAddr == nil || dstAddr == nil {
			_, err := io.WriteString(w, "PROXY UNKNOWN\r\n")
			return err
		}

		family := "TCP4"
		if srcAddr.IP.To4() == nil {
			family = "TCP6"
		}
		_, err := fmt.Fprintf(w, "PROXY %s %s %d %s %d\r\n", family, srcAddr.IP, srcAddr.Port, dstAddr.IP, dstAddr.Port)
		return err
	default:
		return fmt.Errorf("PROXY protocol version %d not supported", dp.ProxyProtocolVersion)
	}
}

// proxyCopy is the function that copies bytes around.
func proxyCopy(errc chan<- error, dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	errc <- err
}

func (dp *DialProxy) keepAlivePeriod() time.Duration {
	if dp.KeepAlivePeriod != 0 {
		return dp.KeepAlivePeriod
	}
	return time.Minute
}

func (dp *DialProxy) dialTimeout() time.Duration {
	if dp.DialTimeout > 0 {
		return dp.DialTimeout
	}
	return 10 * time.Second
}

var defaultDialer = new(net.Dialer)

func (dp *DialProxy) dialContext() func(ctx context.Context, network, address string) (net.Conn, error) {
	if dp.DialContext != nil {
		return dp.DialContext
	}
	return defaultDialer.DialContext
}

func (dp *DialProxy) onDialError() func(src net.Conn, dstDialErr error) {
	if dp.OnDialError != nil {
		return dp.OnDialError
	}
	return func(src net.Conn, dstDialErr error) {
		log.Printf("tcpproxy: for incoming conn %v, error dialing %q: %v", src.RemoteAddr().String(), dp.Addr, dstDialErr)
		src.Close()
	}
}
