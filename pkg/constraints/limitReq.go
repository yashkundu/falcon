package constraints

import (
	"log"
	"net/http"

	"github.com/yashkundu/falcon/pkg/parsing"
	"github.com/yashkundu/falcon/pkg/utils"
)

var limitReqCache *utils.ECache

func getReqCount(key string) int {
	hashkey := utils.Hash(key)
	log.Printf("Hashed Key -> %s", hashkey)
	obj, ok := limitReqCache.Get(hashkey)
	if ok {
		count := obj
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
	key = ip + req.URL.Path
	log.Printf("key -> %s", key)
	count := getReqCount(key)
	if count >= parsing.GetConfig().LimitReq.Frequency {
		return true
	} else {
		return false
	}
}
