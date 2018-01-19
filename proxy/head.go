// Copyright 2017 The margin Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"bufio"
	"bytes"
	"context"
	"net/http"
	"strconv"
	"log"
)

func (p *Proxy) AddHTTPHostRoute(httpHost string, dest Target) {
	p.AddHTTPHostMatchRoute(equals(httpHost), dest)
}

func (p *Proxy) AddHTTPHostMatchRoute(match Matcher, dest Target) {
	p.addRoute(httpHostMatch{match, dest})
}

//no need matcher
type httpHostMatch struct {
	matcher Matcher
	target  Target
}

func (m httpHostMatch) match(br *bufio.Reader) Target {
	if m.matcher(context.TODO(), httpHostHeader(br)) {
		return m.target
	}
	return nil
}

// return the key target in this msg
func tcpKey(br *bufio.Reader) string {
	const maxPeek = 4 << 10
	peekSize := 1
	b, _ := br.Peek(peekSize)
	n := br.Buffered()
	if n > peekSize {
		b, _ = br.Peek(n)
		peekSize = n
	}
	if len(b) > 2 {
		if b[0] != '*' || bytes.Count(b, crlf) < 5 {
			return ""
		}
		fin := 1+bytes.IndexByte(b, '$')
		sin := 1+fin+ bytes.IndexByte(b[fin:], '$')
		tin := sin+bytes.IndexByte(b[sin:], '\r')
		keylen, err := strconv.Atoi(string(b[sin: tin]))
		if err != nil {
			return ""
		}
		return string(b[tin+2 : tin+2+keylen])
	}
	return ""
}

// httpHostHeader returns the HTTP Host header from br without
// consuming any of its bytes. It returns "" if it can't find one.
func httpHostHeader(br *bufio.Reader) string {
	const maxPeek = 4 << 10
	peekSize := 0
	for {
		peekSize++
		log.Printf("peeksize:%d\n",peekSize)
		if peekSize > maxPeek {
			b, _ := br.Peek(br.Buffered())
			log.Printf("--buffered len:%d\n",len(b))
			return httpHostHeaderFromBytes(b)
		}
		b, err := br.Peek(peekSize)
		if n := br.Buffered(); n > peekSize {
			b, _ = br.Peek(n)
			log.Printf("++buffered len:%d\n",n)
			peekSize = n
		}
		if len(b) > 0 {
			if b[0] < 'A' || b[0] > 'Z' {
				// Doesn't look like an HTTP verb
				// (GET, POST, etc).
				return ""
			}
			if bytes.Index(b, crlfcrlf) != -1 || bytes.Index(b, lflf) != -1 {
				req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(b)))
				if err != nil {
					return ""
				}
				if len(req.Header["Host"]) > 1 {
					// TODO(bradfitz): what does
					// ReadRequest do if there are
					// multiple Host headers?
					return ""
				}
				return req.Host
			}
		}
		if err != nil {
			return httpHostHeaderFromBytes(b)
		}
	}
}

var (
	lfHostColon = []byte("\nHost:")
	lfhostColon = []byte("\nhost:")
	crlf        = []byte("\r\n")
	lf          = []byte("\n")
	crlfcrlf    = []byte("\r\n\r\n")
	lflf        = []byte("\n\n")
)

func httpHostHeaderFromBytes(b []byte) string {
	if i := bytes.Index(b, lfHostColon); i != -1 {
		return string(bytes.TrimSpace(untilEOL(b[i+len(lfHostColon):])))
	}
	if i := bytes.Index(b, lfhostColon); i != -1 {
		return string(bytes.TrimSpace(untilEOL(b[i+len(lfhostColon):])))
	}
	return ""
}

// untilEOL returns v, truncated before the first '\n' byte, if any.
// The returned slice may include a '\r' at the end.
func untilEOL(v []byte) []byte {
	if i := bytes.IndexByte(v, '\n'); i != -1 {
		return v[:i]
	}
	return v
}
