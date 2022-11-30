package random

import (
	"math/rand"
	"net/url"

	"github.com/yashkundu/falcon/pkg/balancer"
	"github.com/yashkundu/falcon/pkg/dynamic"
)

// Individual roundrobin balancer for each route
type random struct {
	Servers []*balancer.Server
}

func NewBalancer(urls []*url.URL, varNames []string) *random {
	rr := &random{}

	for i := 0; i < len(urls); i++ {
		rr.AddServer(urls[i], varNames[i])
	}
	return rr
}

func (rr *random) Next() *balancer.Server {
	return rr.Servers[rand.Int()%len(rr.Servers)]
}

func (rr *random) AddServer(url *url.URL, varName string) {
	srv := &balancer.Server{URL: url}
	if varName != "" {
		dynamic.DyServers[varName] = srv
	}
	rr.Servers = append(rr.Servers, srv)
}

func (rr *random) GetServers() []*balancer.Server {
	return rr.Servers
}

func (rr *random) CountServers() int {
	return len(rr.Servers)
}

func (rr *random) SetDead(s *balancer.Server, d bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.IsDead = d
}

func (rr *random) GetIsDead(s *balancer.Server) bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.IsDead
}
