package apiserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/yashkundu/falcon/apiserver/router"
	"github.com/yashkundu/falcon/pkg/parsing"
)

func Run() *http.Server {

	muxRouter := mux.NewRouter()
	router.InitRouter(muxRouter)

	port := parsing.GetConfig().Core.ApiPort
	if port == 0 {
		port = 9900
	}

	srv := &http.Server{
		Handler:      muxRouter,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err.Error())
		}
	}()

	return srv
}
