package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yashkundu/falcon/pkg/proxy"
)

func main() {

	srv := new(proxy.GateServer)
	ss := srv.Run()

	//wait exit
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-quit
	fmt.Println("Start shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, s := range ss {
		if err := s.Shutdown(ctx); err != nil {
			fmt.Printf("Server Shutdown error, signal:%v error:%s\n", sig, err)
		}
	}

	fmt.Printf("Safe exit server, signal:%v\n", sig)
}
