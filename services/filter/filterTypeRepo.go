package filter

import (
	"github.com/omkarp02/pro/db"
)

type FilterTypeRepo struct {
	*db.Database
	collName string
}

func NewRepoFilterType(curDb *db.Database, collName string) *FilterTypeRepo {
	store := &FilterTypeRepo{
		Database: curDb,
		collName: collName,
	}

	store.

	return store
}
