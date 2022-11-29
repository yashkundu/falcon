package wroundrobin

import (
	"net/url"
	"sync/atomic"

	"github.com/yashkundu/falcon/pkg/balancer"
)

// Individual roundrobin balancer for each route
type wroundrobin struct {
	Servers []*balancer.Server
	wts     []uint32
	totalWt uint32
	next    uint32
}

func gcd(a, b uint32) uint32 {
	if a == 0 {
		return b
	}
	return gcd(b%a, a)
}

func (wrr *wroundrobin) reduce() {
	var num uint32 = 1
	for _, s := range wrr.Servers {
		num = gcd(num, s.Wt)
	}
	wrr.totalWt /= num
	for i := 0; i < len(wrr.wts); i++ {
		wrr.wts[i] /= num
	}
}

func NewBalancer(urls []*url.URL) *wroundrobin {
	wrr := &wroundrobin{}
	for i := 0; i < len(urls); i++ {
		wrr.AddServer(urls[i])
	}
	wrr.reduce()
	return wrr
}

func upperBound(wts []uint32, x uint32) int {
	n := len(wts)
	l := 0
	r := n - 1

	for l < r {
		mid := (l + r) / 2
		if wts[mid] <= x {
			l = mid + 1
		} else {
			r = mid
		}
	}
	return l - 1
}

func (wrr *wroundrobin) Next() *balancer.Server {
	n := atomic.AddUint32(&wrr.next, 1)
	ind := upperBound(wrr.wts, (n-1)%wrr.totalWt+uint32(1))
	return wrr.Servers[ind]
}

func (wrr *wroundrobin) AddServer(url *url.URL) {
	curWt := 0
	if curWt < 1 {
		curWt = 1
	}
	wrr.Servers = append(wrr.Servers, &balancer.Server{URL: url})
	wrr.totalWt += uint32(curWt)
	wrr.wts = append(wrr.wts, wrr.wts[len(wrr.wts)-1]+uint32(curWt))
}

func (wrr *wroundrobin) GetServers() []*balancer.Server {
	return wrr.Servers
}

func (wrr *wroundrobin) CountServers() int {
	return len(wrr.Servers)
}

func (wrr *wroundrobin) SetDead(s *balancer.Server, d bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.IsDead = d
}

func (wrr *wroundrobin) GetIsDead(s *balancer.Server) bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.IsDead
}
