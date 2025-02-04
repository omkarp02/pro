package db

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/omkarp02/pro/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	DB     *mongo.Client
	DBName string
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	commandMonitor := &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			// Convert BSON to extended JSON for proper formatting
			formattedCommand, _ := bson.MarshalExtJSON(evt.Command, true, true)
			fmt.Printf("MongoDB Query: %s\n%s\n\n", evt.CommandName, string(formattedCommand))
		},
	}

	uri := cfg.Storage.DBUrl
	clientOptions := options.Client().ApplyURI(uri).SetMonitor(commandMonitor)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	//ping the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Connected to mongodb...")

	return &Database{
		DB:     client,
		DBName: cfg.Storage.DBName,
	}, nil
}
