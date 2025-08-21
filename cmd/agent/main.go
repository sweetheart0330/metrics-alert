package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sweetheart0330/metrics-alert/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	if err := app.RunAgent(ctx); err != nil {
		log.Fatal(err)
	}

	<-ch
	cancel()
}
