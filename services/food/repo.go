package food

import (
	"github.com/omkarp02/pro/db"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repo struct {
	*db.Database
	collName string
}

func NewStore(curDb *db.Database, collName string) *Repo {
	store := &Repo{
		Database: curDb,
		collName: collName,
	}

	return store
}

func (s *Repo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

// func createFood(ctx context.Context, ) error {

// }
