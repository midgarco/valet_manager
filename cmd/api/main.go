package main

import (
	"flag"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/multi"
	"github.com/midgarco/env"
	"github.com/midgarco/valet_manager/internal/services/api"
)

// Flag options that provide values for reading config files and
// connecting to data stores.
var (
	Version = "unset"
	Build   = "unset"

	flagConfigPath  string
	flagEnvironment string
)

func init() {
	flag.StringVar(&flagConfigPath, "config-path", "./config", "path to the config file")
	flag.StringVar(&flagEnvironment, "env", "local", "environment")
}

func main() {
	flag.Parse()

	log.SetHandler(multi.New(
		cli.New(os.Stderr),
	))

	logger := log.WithFields(log.Fields{
		"service":     "api",
		"version":     Version,
		"build":       Build,
		"environment": flagEnvironment,
	})

	// parse the environment file
	err := env.Load("VALET_MGR", flagConfigPath, env.Option{"APP_ENV", flagEnvironment})
	if err != nil {
		logger.Warnf("Could not parse configuration file '%s': %v", flagConfigPath, err)
		return
	}

	api.StartServer(logger)
}
