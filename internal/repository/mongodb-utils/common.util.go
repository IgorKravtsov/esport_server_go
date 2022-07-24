package mongodb_utils

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertToObjID(possibleID interface{}) (primitive.ObjectID, error) {
	oid, ok := possibleID.(primitive.ObjectID)
	if !ok {
		return primitive.ObjectID{}, errors.New("failed to convert objectID to HEX")
	}
	return oid, nil
}
