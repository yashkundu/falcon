package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yashkundu/falcon/apiserver"
	"github.com/yashkundu/falcon/pkg/parsing"
	"github.com/yashkundu/falcon/pkg/proxy"
)

func main() {

	p := new(proxy.GateServer)
	servers := p.Run()

	serverPort := parsing.GetConfig().Core.Listen
	if serverPort == 0 {
		serverPort = 80
	}

	log.Printf("Falcon running in the background on port %d ... ", serverPort)

	apiSrv := apiserver.Run()
	apiPort := parsing.GetConfig().Core.ApiPort
	if apiPort == 0 {
		apiPort = 9900
	}
	log.Printf("Falcon Api server running on port %d ... ", apiPort)

	servers = append(servers, apiSrv)

	//wait exit
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-quit
	fmt.Println("Start shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, s := range servers {
		if err := s.Shutdown(ctx); err != nil {
			fmt.Printf("Server Shutdown error, signal:%v error:%s\n", sig, err)
		}
	}

	fmt.Printf("Safe exit server, signal:%v\n", sig)
}
