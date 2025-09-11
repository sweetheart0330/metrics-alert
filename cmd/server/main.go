package main

import (
	"log"

	"github.com/sweetheart0330/metrics-alert/internal/app"
)

func main() {
	if err := app.RunServer(); err != nil {
		log.Printf("server error, err: %v, exit", err)
	}
}
