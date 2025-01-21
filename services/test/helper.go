package test

import "go.mongodb.org/mongo-driver/v2/bson"

func helper(data string) error {

	_, err := bson.ObjectIDFromHex(data)
	if err != nil {
		return err
	}

	return nil

}
