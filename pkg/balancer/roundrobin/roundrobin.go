package roundrobin

import (
	"net/url"
	"sync/atomic"

	"github.com/yashkundu/falcon/pkg/balancer"
)

// Individual roundrobin balancer for each route
type roundrobin struct {
	Servers []balancer.Server
	next    uint32
}

func NewBalancer(urls []*url.URL) (*roundrobin, error) {
	rr := &roundrobin{}

	for i := 0; i < len(urls); i++ {
		rr.AddServer(urls[i])
	}
	return rr, nil
}

func (rr *roundrobin) Next() *balancer.Server {
	n := atomic.AddUint32(&rr.next, 1)
	return &rr.Servers[(int(n)-1)%len(rr.Servers)]
}

func (rr *roundrobin) AddServer(url *url.URL) {
	rr.Servers = append(rr.Servers, balancer.Server{URL: url})
}

func (rr *roundrobin) GetServers() *[]balancer.Server {
	return &rr.Servers
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
