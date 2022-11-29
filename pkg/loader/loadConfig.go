package loader

import (
	"log"
	"net/url"
	"time"

	"github.com/yashkundu/falcon/pkg/balancer"
	"github.com/yashkundu/falcon/pkg/balancer/random"
	"github.com/yashkundu/falcon/pkg/balancer/roundrobin"
	"github.com/yashkundu/falcon/pkg/parsing"
	"github.com/yashkundu/falcon/pkg/utils"
)

type BalancerType int

const (
	RandomBalancer      BalancerType = 1
	RoundRobinBalancer  BalancerType = 2
	WRoundRobinBalancer BalancerType = 3
)

type MatchType int

const (
	ExactMatch  MatchType = 1
	PrefixMatch MatchType = 2
	RegexMatch  MatchType = 3
)

type RouteInfo struct {
	Endpoint      string
	Match         MatchType
	TargetServer  *balancer.Server
	IsMultiTarget bool
	Balancer      balancer.Balancer
}

var (
	LimitReqCache *utils.ECache
	Routes        []*RouteInfo
)

func loadConfig() error {
	config := parsing.GetConfig()

	if config.LimitReq.Enable {
		LimitReqCache = utils.NewECache(time.Millisecond * time.Duration(config.LimitReq.Interval))
	}

	proxyConfig := config.Proxy
	for _, routeConfig := range proxyConfig.Routes {
		routeInfo := getRouteInfo(&routeConfig)
		Routes = append(Routes, routeInfo)
	}

	return nil
}

func getRouteInfo(routeConfig *parsing.Route) *RouteInfo {

	routeInfo := new(RouteInfo)

	if len(routeConfig.Backends) == 0 {
		panic(routeConfig.Endpoint + "Empty target")
	}

	routeInfo.Endpoint = routeConfig.Endpoint
	routeInfo.Match = MatchType(routeConfig.Match)

	if len(routeConfig.Backends) == 1 {
		routeInfo.IsMultiTarget = false
		url, err := url.Parse(routeConfig.Backends[0].Url)
		if err != nil {
			log.Fatal(err)
		}
		routeInfo.TargetServer = &balancer.Server{URL: url}
		return routeInfo
	}

	routeInfo.IsMultiTarget = true
	urls := make([]*url.URL, 0)
	for _, b := range routeConfig.Backends {
		curUrl, err := url.Parse(b.Url)
		if err != nil {
			log.Fatal(err)
		}
		urls = append(urls, curUrl)
	}

	switch BalancerType(routeConfig.Balancer) {
	case RandomBalancer:
		routeInfo.Balancer = random.NewBalancer(urls)
	case RoundRobinBalancer:
		routeInfo.Balancer = roundrobin.NewBalancer(urls)
	case WRoundRobinBalancer:
		routeInfo.Balancer = roundrobin.NewBalancer(urls)
	}

	return routeInfo
}

func (routeInfo *RouteInfo) GetTargetServer() *balancer.Server {
	if !routeInfo.IsMultiTarget {
		return routeInfo.TargetServer
	}
	return routeInfo.Balancer.Next()
}

func Init() error {
	loadConfig()
	return nil
}
