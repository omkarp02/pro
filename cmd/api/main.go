package main

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/omkarp02/pro/api"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
)

func main() {

	cfg := config.MustLoad()
	config.SetUpLogger()

	// mongoClient, err := db.ConnectToDB(cfg)
	// if err != nil {
	// 	log.Fatal(err)
	// }

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

	slog.Info("here it reached")

	server := api.NewAPIServer(cfg.Addr, DB, cfg)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
