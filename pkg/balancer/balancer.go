package balancer

import (
	"net/url"
	"sync"
)

type Balancer interface {
	Next() *Server
	AddServer(url *url.URL)
	GetServers() *[]Server
	CountServers() int
}

type Server struct {
	URL    *url.URL
	IsDead bool
	Wt     uint32
	Mu     sync.RWMutex
}
