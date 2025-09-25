package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/sweetheart0330/metrics-alert/internal/app"
)

func main() {
	parentCtx := context.Background()
	ctx, stop := signal.NotifyContext(parentCtx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		if err := app.RunAgent(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	stop()

	fmt.Println("Application stopped")
}
