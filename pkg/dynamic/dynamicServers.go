package dynamic

import (
	"github.com/yashkundu/falcon/pkg/balancer"
)

var (
	DyServers map[string]*balancer.Server
)

func init() {
	DyServers = make(map[string]*balancer.Server, 0)
}
