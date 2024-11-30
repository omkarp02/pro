package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/omkarp02/pro/api"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
)

func main() {

	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags
		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	cfg := config.MustLoad(configPath)
	config.SetUpLogger()

	DB, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = DB.DB.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	server := api.NewAPIServer(cfg.Addr, DB, cfg)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
