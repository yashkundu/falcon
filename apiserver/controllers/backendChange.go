package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/yashkundu/falcon/apiserver/common"
	"github.com/yashkundu/falcon/apiserver/resultcodes"
	"github.com/yashkundu/falcon/pkg/dynamic"
)

type Change struct {
	VarName string `json:"varName"`
	Url     string `json:"url"`
}

func BackendChange(w http.ResponseWriter, r *http.Request) {

	var change Change
	var success bool

	body, err1 := ioutil.ReadAll(r.Body)

	if err1 != nil {
		fmt.Println(err1)
		success = false
	} else {
		err := json.Unmarshal(body, &change)
		if err != nil {
			fmt.Println(err)
			success = false
		} else {
			srv, ok := dynamic.DyServers[change.VarName]
			if !ok {
				success = false
			} else {
				parsedUrl, err := url.Parse(change.Url)
				if err != nil {
					fmt.Println(err)
					success = false
				} else {
					srv.Mu.Lock()
					srv.URL = parsedUrl
					srv.Mu.Unlock()
					success = true
				}
			}
		}
	}

	err2 := common.Report(w, resultcodes.SUCCESS, "success", common.H{"success": success})
	if err2 != nil {
		fmt.Println(err2)
	}
}
