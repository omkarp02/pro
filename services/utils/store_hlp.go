package services

import (
	"errors"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetCurrentTimestamps() Timestamps {
	return Timestamps{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func CreateBsonFromKeyValuePair(fields ...interface{}) (bson.D, error) {
	if len(fields)%2 != 0 {
		return nil, errors.New("invalid number of arguments, must be in key-value pairs")
	}

	var filter bson.D
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			return nil, errors.New("key must be string")
		}
		value := fields[i+1]
		filter = append(filter, bson.E{Key: key, Value: value})
	}

	return filter, nil
}

func CreateBsonFromStruct(input interface{}) bson.D {
	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)

	//check if it a struct pointer then derefrence and get the actual thing
	// Ensure it's a struct
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic("input must be a struct or a pointer to a struct")
	}

	updates := bson.D{}

	// Iterate over the fields
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Get the field name and value
		fieldName := field.Name
		fieldValue := value.Interface()

		updates = append(updates, bson.E{Key: fieldName, Value: fieldValue})
	}

	return updates
}
