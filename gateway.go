package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type AcceptSerevr struct{}

func (s *AcceptSerevr) runProxy(port string) *http.Server {
	ln, err := net.Listen("tcp", ":8009")
	if err != nil {
		panic(err)
	}
	p := &Proxy{}
	srv := &http.Server{Addr: ":" + port, Handler: p}
	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			panic(err.Error())
		}
	}()
	log.Print("proxy server listen on port", port)
	return srv
}

type Proxy struct{}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ip string
	ips := parseIpAddr(r)
	if len(ips) == 0 {
		log.Print("ips", ips)
		return
	}
	ip = ips[0]
	log.Print("accept request from ip", ip)
	if err := p.blackList(ip); err != nil {
		log.Print("error", err)
	}

	if err := p.whiteList(ip); err != nil {
		log.Print("error", err)
	}
	//转发
	p.doProxy(w, r)
}
func (p *Proxy) blackList(ip string) error {
	return nil
}
func (p *Proxy) whiteList(ip string) error {
	return nil
}
func (p *Proxy) doProxy(w http.ResponseWriter, r *http.Request) {
	route := p.getRoute(r)
	target, err := url.Parse(route)
	if err != nil {
		log.Print(" get route err ", err)
		return
	}
	proxy := newReverseProxy(target)

	//todo 基于响应时间做负载优化
	in := time.Now()
	proxy.ServeHTTP(w, r)
	log.Print("time cost", time.Now().Sub(in).Seconds(), "s")
}
func (p *Proxy) getRoute(req *http.Request) string {
	return HostInfo.GetTarget(req)
}
func newReverseProxy(target *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}
		req.Header["X-Real-Ip"] = parseIpAddr(req)
		log.Print("X-Real-Ip=", req.Header["X-Real-Ip"])
	}
	return &httputil.ReverseProxy{Director: director}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
func parseIpAddr(req *http.Request) []string {
	return []string{strings.Split(req.RemoteAddr, ":")[0]}
}
