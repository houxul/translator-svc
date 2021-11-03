package main

import (
	"context"
	"fmt"
	"log"
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
		Handler: router(),
	}
	grace.Register(func() {
		s.Shutdown(context.Background())
	})
	grace.Run(s.ListenAndServe)
}

func router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/translator", translatorHandler)
	mux.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("./resources"))))
	return mux
}

func pingHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "ok")
}

func translatorHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		translator.Translate(w, req)
	default:
		fmt.Fprintf(w, "%s Method not found", req.Method)
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
		log.Printf("receive signal(%s)\n", sig.String())
	case err := <-quit:
		log.Printf("exit with error(%v)\n", err)
	}

	for _, cleanup := range g.cleanups {
		cleanup()
	}
}
