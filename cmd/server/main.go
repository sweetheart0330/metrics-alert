package main

import (
	"context"
	"fmt"
	"github.com/sweetheart0330/metrics-alert/internal/app"
	"log"
)

func main() {
	ctx := context.Background()
	if err := app.RunServer(ctx); err != nil {
		log.Printf("server error, err: %v, exit", err)
	}

	fmt.Println("Server stopped")
}
