package main

import (
	"context"
	"flag"

	"delivery/internal/app"

	"github.com/labstack/gommon/log"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "../../deploy/env/.env.local", "path to config file")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	application, err := app.NewApp(ctx, configPath)
	if err != nil {
		log.Fatalf("failed to create application: %v", err)
	}

	err = application.Run()
	if err != nil {
		log.Fatalf("failed to run application: %v", err)
	}
}
