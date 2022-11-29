package controllers

import (
	"fmt"
	"net/http"

	"github.com/yashkundu/falcon/apiserver/common"
	"github.com/yashkundu/falcon/apiserver/resultcodes"
	"github.com/yashkundu/falcon/pkg/constraints/status"
)

func StatusReqCount(w http.ResponseWriter, r *http.Request) {
	err := common.Report(w, resultcodes.SUCCESS, "success", common.H{"count": status.Instance().GetReqCount()})
	if err != nil {
		fmt.Println(err)
	}
}
