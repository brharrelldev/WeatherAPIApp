package main

import (
	"flag"
	"github.com/brharrelldev/weatherAPI/service"
	"github.com/common-nighthawk/go-figure"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
)

var (
	configPath string
	apiKey     string
)

func main() {

	banner := figure.NewFigure("WEATHER API", "", true)
	banner.Print()

	flag.StringVar(&configPath, "config", "config.yaml", "path to config file")
	flag.StringVar(&apiKey, "apikey", "", "api key")
	flag.Parse()

	conf, err := parseConfig(apiKey, configPath)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewProduction()

	s, err := service.NewService(conf, logger)
	if err != nil {
		log.Fatalf("error creating service: %v", err)
	}

	errChan := make(chan error)
	sigChan := make(chan os.Signal, 1)
	go func() {
		if err := s.Start(apiKey); err != nil {
			errChan <- err
			return
		}
	}()

	signal.Notify(sigChan, os.Interrupt)

	select {
	case err = <-errChan:
		if err != nil {
			panic(err)
		}
	case <-sigChan:
		log.Println("graceful shutdown received")
		os.Exit(0)
	}
}
