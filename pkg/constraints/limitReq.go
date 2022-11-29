package constraints

import (
	"net/http"

	"github.com/yashkundu/falcon/pkg/parsing"
	"github.com/yashkundu/falcon/pkg/utils"
)

var limitReqCache *utils.ECache

func getReqCount(key string) int {
	hashkey := utils.Hash(key)
	obj, ok := limitReqCache.Get(hashkey)
	if ok {
		count := obj.(int)
		//  wrong here
		limitReqCache.Set(hashkey, count+1)
		return count
	} else {
		limitReqCache.Set(hashkey, 1)
	}
	return 0
}

func ExceededLimitReq(ip string, req *http.Request) bool {
	var key string
	if parsing.GetConfig().LimitReq.Mode == 0 {
		key = ip + req.Host + req.URL.Path
	} else {
		key = ip + req.RequestURI
	}
	count := getReqCount(key)
	if count >= parsing.GetConfig().LimitReq.Frequency {
		return true
	} else {
		return false
	}
}
