package roundrobin

import (
	"net/url"
	"sync/atomic"

	"github.com/yashkundu/falcon/pkg/balancer"
	"github.com/yashkundu/falcon/pkg/dynamic"
)

// Individual roundrobin balancer for each route
type roundrobin struct {
	Servers []*balancer.Server
	next    uint32
}

func NewBalancer(urls []*url.URL, varNames []string) *roundrobin {
	rr := &roundrobin{}

	for i := 0; i < len(urls); i++ {
		rr.AddServer(urls[i], varNames[i])
	}
	return rr
}

func (rr *roundrobin) Next() *balancer.Server {
	n := atomic.AddUint32(&rr.next, 1)
	return rr.Servers[(int(n)-1)%len(rr.Servers)]
}

func (rr *roundrobin) AddServer(url *url.URL, varName string) {
	srv := &balancer.Server{URL: url}
	if varName != "" {
		dynamic.DyServers[varName] = srv
	}
	rr.Servers = append(rr.Servers, srv)
}

func (rr *roundrobin) GetServers() []*balancer.Server {
	return rr.Servers
}

func (rr *roundrobin) CountServers() int {
	return len(rr.Servers)
}

func (rr *roundrobin) SetDead(s *balancer.Server, d bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.IsDead = d
}

func (rr *roundrobin) GetIsDead(s *balancer.Server) bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.IsDead
}
