package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/joho/godotenv"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/InVisionApp/kit-overwatch/api"
	"github.com/InVisionApp/kit-overwatch/config"
	"github.com/InVisionApp/kit-overwatch/deps"
	"github.com/InVisionApp/kit-overwatch/watcher"
)

var (
	version       = "No version specified"
	flushInterval = time.Duration(100 * time.Millisecond)

	envFile = kingpin.Flag("envfile", "Specify a different dotenv file to use for loading env vars").Short('f').Default(".env").String()
)

func init() {
	// Parse CLI stuff
	kingpin.Version(version)
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.Parse()

	if err := godotenv.Load(*envFile); err != nil {
		log.Warningf("Unable to load dotenv file '%v': %v", *envFile, err.Error())
	}
}

func main() {
	cfg := config.New()

	if err := cfg.LoadEnvVars(); err != nil {
		log.Fatalf("Configuration error: %v", err.Error())
	}

	// Show debug logs if debug mode enabled
	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug mode enabled")
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Log the notification level
	log.Infof("Notification level set to: %s", cfg.NotificationLevel)

	// create a statsd client
	statsdClient, err := statsd.NewBufferedClient(cfg.StatsDAddress, cfg.StatsDPrefix, flushInterval, 0)
	if err != nil {
		log.Fatalf("Unable to instantiate statsd client: %v", err.Error())
	}

	// For dependency injection
	d := &deps.Dependencies{
		StatsD: statsdClient,
	}

	// Start the watcher
	w := watcher.New(cfg)
	go w.Watch()

	// Start the API server
	api := api.New(cfg, d, version)
	log.Fatal(api.Run())
}
