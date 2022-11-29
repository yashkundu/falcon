package status

import (
	"sync"
)

var status *Status

type Status struct {
	ReqCount   int
	ReqCountMu *sync.RWMutex
}

// can also use atomic Adduint32 or something else

func Init() error {
	status = &Status{
		ReqCount:   0,
		ReqCountMu: new(sync.RWMutex),
	}
	return nil
}

func Instance() *Status {
	return status
}

func (sta *Status) AddReqCount() {
	sta.ReqCountMu.Lock()
	sta.ReqCount++
	sta.ReqCountMu.Unlock()
}

func (sta *Status) SubReqCount() {
	sta.ReqCountMu.Lock()
	sta.ReqCount--
	sta.ReqCountMu.Unlock()
}

func (sta *Status) GetReqCount() int {
	sta.ReqCountMu.RLock()
	defer sta.ReqCountMu.RUnlock()
	return sta.ReqCount
}
