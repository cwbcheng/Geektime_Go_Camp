package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var errCh = make(chan error, 1)

		srv := http.Server{
			Addr: ":18081",
			Handler: http.DefaultServeMux,
		}
		go func() {
			defer func() {
				fmt.Println("Listening goroutine quit.")
			}()

			errCh <- srv.ListenAndServe()
		}()

		select {
		case err := <- errCh:
			return err
		case <- ctx.Done():
			fmt.Println("Receive cancel.")
			return srv.Close()
		}
	})
	
	signals := make(chan os.Signal, 1)
	shutdownArray :=
		[]os.Signal {syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT}
	signal.Notify(signals, shutdownArray...)

	defer func() {
		cancel()
		close(signals)
	}()

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