package db

import (
	"context"
)

type MongoTransactionManager struct {
	db *Database
}

type InTxn func(sessCtx context.Context) (interface{}, error)

// TransactionManager defines the interface for transaction execution
type TransactionManager interface {
	RunInTxn(ctx context.Context, fn InTxn) (interface{}, error)
}

// NewMongoTransactionManager creates a new MongoTransactionManager
func NewMongoTransactionManager(db *Database) *MongoTransactionManager {
	return &MongoTransactionManager{db: db}
}

func (tm *MongoTransactionManager) RunInTxn(ctx context.Context, fn InTxn) (interface{}, error) {
	// Start a new session
	session, err := tm.db.DB.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	// Run the transaction
	result, err := session.WithTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
		return fn(sessCtx)
	})
	return result, err
}
