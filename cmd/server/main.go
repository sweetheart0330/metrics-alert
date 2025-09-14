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

	if err := app.RunServer(ctx); err != nil {
		log.Printf("server error, err: %v, exit", err)
	}

	<-ctx.Done()
	stop()

	fmt.Println("Server stopped")
}
