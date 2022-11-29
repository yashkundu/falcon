package random

import (
	"math/rand"
	"net/url"

	"github.com/yashkundu/falcon/pkg/balancer"
)

// Individual roundrobin balancer for each route
type random struct {
	Servers []*balancer.Server
}

func NewBalancer(urls []*url.URL) *random {
	rr := &random{}

	for i := 0; i < len(urls); i++ {
		rr.AddServer(urls[i])
	}
	return rr
}

func (rr *random) Next() *balancer.Server {
	return rr.Servers[rand.Int()%len(rr.Servers)]
}

func (rr *random) AddServer(url *url.URL) {
	rr.Servers = append(rr.Servers, &balancer.Server{URL: url})
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
