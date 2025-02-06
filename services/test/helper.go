package test

import "go.mongodb.org/mongo-driver/v2/bson"

func helper(data string) error {

	_, err := bson.ObjectIDFromHex(data)
	if err != nil {
		return err
	}

	return nil

	// objectIds, err := store.SliceOfHexToObjectID(Id)

	// userId := objectIds[0]
	// addressId := objectIds[1]

}

// result, err := s.txn.RunInTxn(ctx, func(sessCtx context.Context) (interface{}, error) {})
