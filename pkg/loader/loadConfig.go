package loader

import (
	"log"
	"net/url"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/yashkundu/falcon/pkg/balancer"
	"github.com/yashkundu/falcon/pkg/balancer/random"
	"github.com/yashkundu/falcon/pkg/balancer/roundrobin"
	"github.com/yashkundu/falcon/pkg/dynamic"
	"github.com/yashkundu/falcon/pkg/parsing"
	"github.com/yashkundu/falcon/pkg/utils"
)

type BalancerType int

const (
	RoundRobinBalancer  BalancerType = 0
	RandomBalancer      BalancerType = 1
	WRoundRobinBalancer BalancerType = 2
)

type MatchType int

const (
	ExactMatch  MatchType = 0
	PrefixMatch MatchType = 1
	RegexMatch  MatchType = 2
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
		if routeConfig.Backends[0].VarName != "" {
			dynamic.DyServers[routeConfig.Backends[0].VarName] = routeInfo.TargetServer
		}
		return routeInfo
	}

	routeInfo.IsMultiTarget = true
	urls := make([]*url.URL, 0)
	varNames := make([]string, 0)
	for _, b := range routeConfig.Backends {
		curUrl, err := url.Parse(b.Url)
		if err != nil {
			log.Fatal(err)
		}
		urls = append(urls, curUrl)
		varNames = append(varNames, b.VarName)
	}

	switch BalancerType(routeConfig.Balancer) {
	case RandomBalancer:
		routeInfo.Balancer = random.NewBalancer(urls, varNames)
	case RoundRobinBalancer:
		routeInfo.Balancer = roundrobin.NewBalancer(urls, varNames)
	case WRoundRobinBalancer:
		routeInfo.Balancer = roundrobin.NewBalancer(urls, varNames)
	}

	return routeInfo
}

func (routeInfo *RouteInfo) GetTargetServer() *balancer.Server {
	if !routeInfo.IsMultiTarget {
		return routeInfo.TargetServer
	}
	return routeInfo.Balancer.Next()
}

func init() {
	loadConfig()
	log.Println("Loaded Configs : ")
	spew.Dump(Routes)
}
