package main

import (
	"context"
	"log"
	"runtime/pprof"
)

func main() {
	chat := NewChat()
	hub := NewHub()
	//go hub.run()
	go func() {
		labels := pprof.Labels("func", "run")
		pprof.Do(context.Background(), labels, func(_ context.Context) {
			hub.run()
		})
	}()
	chat.rooms["Общая"] = hub

	err := chat.Router().Run(":81")
	if err != nil {
		log.Fatalf("FATAL: Error starting server: %s\n", err)
	}
}
