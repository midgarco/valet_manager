package main

import (
	"flag"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/multi"
	"github.com/midgarco/valet_manager/internal/services/http"
	"github.com/midgarco/valet_manager/pkg/config"
)

var (
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
		"environment": flagEnvironment,
	})

	// parse the environment file
	err := config.LoadEnv(flagConfigPath, config.Option{"APP_ENV", flagEnvironment})
	if err != nil {
		logger.Warnf("Could not parse configuration file '%s': %v", flagConfigPath, err)
		return
	}

	http.StartServer(logger)
}
