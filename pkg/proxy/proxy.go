package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/yashkundu/falcon/pkg/balancer"
	"github.com/yashkundu/falcon/pkg/constraints"
	"github.com/yashkundu/falcon/pkg/constraints/status"
	"github.com/yashkundu/falcon/pkg/loader"
	"github.com/yashkundu/falcon/pkg/parsing"
	"golang.org/x/net/netutil"
)

type Proxy struct {
}

func (p *Proxy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	spew.Dump(req.URL.Path, req.URL.Host)

	if parsing.GetConfig().LimitReq.Enable {
		ip := strings.Split(req.RemoteAddr, ":")[0]
		if constraints.ExceededLimitReq(ip, req) {
			return
		}
	}

	if parsing.GetConfig().Core.EnableServerStats {
		status.Instance().AddReqCount()
		defer status.Instance().SubReqCount()
	}

	route := new(loader.RouteInfo)
	existRoute := false

	for _, r := range loader.Routes {
		spew.Dump("req url path", req.URL.Path)
		spew.Dump("route endpoint", r.Endpoint)
		switch r.Match {
		case loader.ExactMatch:
			if req.URL.Path == r.Endpoint {
				route = r
				existRoute = true
			}
		case loader.PrefixMatch:
			if strings.HasPrefix(req.URL.Path, r.Endpoint) {
				route = r
				existRoute = true
			}
		case loader.RegexMatch:
			if match, _ := regexp.MatchString(r.Endpoint, req.URL.Path); match {
				route = r
				existRoute = true
			}
		}
	}

	spew.Dump(existRoute, route)

	if existRoute {
		proxy := newHostReverseProxy(route.GetTargetServer())
		proxy.ServeHTTP(res, req)
	}
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

func newHostReverseProxy(server *balancer.Server) *httputil.ReverseProxy {
	spew.Dump(server)
	director := func(req *http.Request) {
		targetQuery := server.URL.RawQuery
		req.URL.Scheme = server.URL.Scheme
		req.URL.Host = server.URL.Host
		req.URL.Path = singleJoiningSlash(server.URL.Path, req.URL.Path)
		req.Host = server.URL.Host

		spew.Dump(req.URL.Scheme, req.URL.Host, req.URL.Path)

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}
		xForwarded := req.Header["X-Forwarded-For"]
		xForwarded = append(xForwarded, req.RemoteAddr)
		req.Header["X-Forwarded-For"] = xForwarded
	}
	return &httputil.ReverseProxy{Director: director}
}

type GateServer struct{}

// proxy80 -- http , proxy443 -- https (Not implemented now), proxyws -- websocket (Not implemented untill now)
func (s *GateServer) proxy80() *http.Server {
	port := parsing.GetConfig().Core.Listen
	if port == 0 {
		port = 80
	}
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	if parsing.GetConfig().Core.LimitMaxConn > 0 {
		ln = netutil.LimitListener(ln, parsing.GetConfig().Core.LimitMaxConn)
	}

	p := &Proxy{}
	srv := &http.Server{Addr: ":80", Handler: p}
	if parsing.GetConfig().Core.ReadTimeout > 0 {
		srv.ReadTimeout = time.Duration(parsing.GetConfig().Core.ReadTimeout) * time.Second
	}
	if parsing.GetConfig().Core.WriteTimeout > 0 {
		srv.WriteTimeout = time.Duration(parsing.GetConfig().Core.WriteTimeout) * time.Second
	}
	if parsing.GetConfig().Core.IdleTimeout > 0 {
		srv.IdleTimeout = time.Duration(parsing.GetConfig().Core.IdleTimeout) * time.Second
	}

	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			panic(err.Error())
		}
	}()
	return srv
}

func (s *GateServer) Run() []*http.Server {
	servers := make([]*http.Server, 0)
	p80 := s.proxy80()
	servers = append(servers, p80)
	return servers
}
