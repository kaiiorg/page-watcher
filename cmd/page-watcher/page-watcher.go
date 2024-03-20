package main

import (
	"flag"

	"github.com/kaiiorg/page-watcher/pkg/config"
	"github.com/kaiiorg/page-watcher/pkg/util"
	"github.com/kaiiorg/page-watcher/pkg/watcher"

	"github.com/rs/zerolog/log"
)

const (
	applicationName        = "page-watcher"
	applicationDescription = "Watches a web page for changes. See https://github.com/kaiiorg/page-watcher"
)

var (
	logLevel   = flag.String("log-level", "info", "Zerolog log level to use; trace, debug, info, warn, error, panic, etc")
	configPath = flag.String("config", "./configs/config.hcl", "path to HCL config file")
)

func main() {
	flag.Parse()

	util.ConfigureLogging(*logLevel, applicationName, applicationDescription)
	conf, err := config.LoadFromFile(*configPath)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("path", *configPath).
			Msg("Failed to load config file")
	}

	log.Info().Int("pages", len(conf.Pages)).Msg("Loaded configuration")

	watcher, err := watcher.New(conf)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to create watcher")
	}

	watcher.Watch()
	util.WaitForInterrupt()
	watcher.Close()
}
