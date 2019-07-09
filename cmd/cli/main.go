package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/multi"
	"github.com/midgarco/env"
	"github.com/midgarco/valet_manager/internal/manager"
	"github.com/midgarco/valet_manager/pkg/valet"
)

// Flag options that provide values for reading config files and
// connecting to data stores.
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
		"service":     "cli",
		"environment": flagEnvironment,
	})

	// parse the environment file
	err := env.Load("VALET_MGR", flagConfigPath, env.Option{"APP_ENV", flagEnvironment})
	if err != nil {
		logger.Warnf("Could not parse configuration file '%s': %v", flagConfigPath, err)
		return
	}

	// redis connection
	conn := &manager.Connection{}
	if err := conn.RedisConnection(); err != nil {
		logger.WithError(err).Fatalf("error occured connecting to redis")
	}

	// database connection
	if err := conn.DBConnection(); err != nil {
		logger.WithError(err).Fatalf("error occured connecting to database")
	}
	defer conn.DB.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Valet Manager CLI Tool")
	fmt.Println("--------------------------------")

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')

		// trim new lines
		text = strings.Trim(text, "\n")

		switch text {
		case "create super user":
			fmt.Print("creating super admin user in database: ")
			if err := createSuperUser(conn); err != nil {
				fmt.Printf("failed creating user: %v\n", err)
			}
			fmt.Println("done")
		}
	}
}

func createSuperUser(conn *manager.Connection) error {
	user := valet.User{
		FirstName: "Jeff",
		LastName:  "Dupont",
		Email:     "jeff.dupont@gmail.com",
		Password:  env.Get("SUPER_PASS"),
		Address: valet.Address{
			Line1:   "123 Main St",
			City:    "Anycity",
			State:   "CA",
			Zipcode: "00001",
		},
		PhoneNumbers: []valet.PhoneNumber{
			valet.PhoneNumber{Type: "home", Value: "222 123-4567"},
			valet.PhoneNumber{Type: "work", Value: "333 456-7890"},
			valet.PhoneNumber{Type: "mobile", Value: "444 567-8901"},
		},
		Admin: true,
	}

	return user.Create(conn.DB)
}
