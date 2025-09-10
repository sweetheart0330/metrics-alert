package main

import (
	"log"

	"github.com/sweetheart0330/metrics-alert/internal/app"
)

func main() {
	if err := app.RunServer(); err != nil {
		log.Fatal(err)
	}
}
