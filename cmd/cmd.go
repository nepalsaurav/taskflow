package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/nepalsaurav/taskflow/api"
	"github.com/nepalsaurav/taskflow/models"
)

type ServeConfig struct {
	Address string
}

func initServe(args []string) {
	config := &ServeConfig{
		Address: "127.0.0.1:8080",
	}
	flags := flag.NewFlagSet("serve", flag.ExitOnError)
	flags.StringVar(&config.Address, "address", "127.0.0.1:8080", "serve address")
	flags.Parse(args)

	// migrate
	models.Migrate()
	// serve
	api.Serve(config.Address)
}

func InitCMD() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("expected command like serve, create_api, migrate")
	}

	switch os.Args[1] {
	case "serve":
		{
			initServe(os.Args[2:])
		}
	default:
		return fmt.Errorf("flag not implemented")
	}
	return nil
}
