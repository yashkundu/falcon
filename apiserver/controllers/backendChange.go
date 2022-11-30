package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/yashkundu/falcon/apiserver/common"
	"github.com/yashkundu/falcon/apiserver/resultcodes"
	"github.com/yashkundu/falcon/pkg/dynamic"
)

type Change struct {
	varName string
	url     string
}

func BackendChange(w http.ResponseWriter, r *http.Request) {

	var change Change
	var success bool

	err1 := json.NewDecoder(r.Body).Decode(&change)

	if err1 != nil {
		fmt.Println(err1)
		success = false
	} else {
		srv, ok := dynamic.DyServers[change.varName]
		if !ok {
			success = false
		} else {
			parsedUrl, err := url.Parse(change.url)
			if err != nil {
				fmt.Println(err)
				success = false
			} else {
				srv.URL = parsedUrl
				success = true
			}
		}
	}

	err2 := common.Report(w, resultcodes.SUCCESS, "success", common.H{"success": success})
	if err2 != nil {
		fmt.Println(err2)
	}
}
