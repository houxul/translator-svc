package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"translator/provider"
	"translator/translator"
)

func main() {
	var grace grace
	grace.Register(provider.Engine.Close)

	s := &http.Server{
		Addr:    ":8090",
		Handler: http.HandlerFunc(router),
	}
	grace.Register(func() {
		s.Shutdown(context.Background())
	})
	grace.Run(s.ListenAndServe)
}

func router(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/ping" {
		fmt.Fprintln(w, "ok")
		return
	}

	switch req.Method {
	case http.MethodPost:
		switch req.URL.Path {
		case "/translator":
			translator.Translate(w, req)
		default:
			fmt.Fprintf(w, "%s Path not found", req.URL.Path)
			return
		}
	default:
		fmt.Fprintf(w, "%s Method not found", req.Method)
		return
	}
}

type grace struct {
	cleanups []func()
}

func (g *grace) Register(f func()) {
	g.cleanups = append(g.cleanups, f)
}

func (g *grace) Run(f func() error) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	quit := make(chan error)
	go func() {
		var err error
		defer func() {
			if e := recover(); e != nil {
				quit <- fmt.Errorf("painc:%v", e)
			} else {
				quit <- err
			}
		}()

		err = f()
	}()

	select {
	case sig := <-signals:
		fmt.Printf("receive signal(%s)\n", sig.String())
	case err := <-quit:
		fmt.Printf("exit with error(%v)\n", err)
	}

	for _, cleanup := range g.cleanups {
		cleanup()
	}
}
