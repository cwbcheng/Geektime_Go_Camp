package main

import (
	pb "Geektime_Go_Camp/第四周/api"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	//"Geektime_Go_Camp/第四周/internal/service"
)

const (
	port = ":50051"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	//注册服务，监听端口
	g.Go(func() error {
		var errCh = make(chan error, 1)
		service := InitUserService()
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			return err
		}
		s := grpc.NewServer()
		pb.RegisterUserServiceServer(s, service)
		log.Printf("server listening at %v", lis.Addr())

		go func() {
			if err := s.Serve(lis); err != nil {
				log.Fatalf("failed to serve: %v", err)
				errCh <- err
			}
		}()

		select {
		case err := <- errCh:
			return err
		case <- ctx.Done():
			fmt.Println("Receive cancel.")
			s.Stop()
			return nil
		}
	})

	signals := make(chan os.Signal, 1)
	shutdownArray :=
		[]os.Signal {syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT}
	signal.Notify(signals, shutdownArray...)

	//监听信号
	g.Go(func() error {
		for {
			sig := <- signals
			for _, item := range shutdownArray {
				if sig == item {
					fmt.Println("Receive quit signal.")
					cancel()
					go func() {
						select {
						case <- time.After(5 * time.Second):
							fmt.Println("time out")
							os.Exit(0)
						}
					}()
					return nil
				}
			}
		}
	})

	err := g.Wait()
	if err != nil {
		log.Fatalln(err)
		return
	}
}