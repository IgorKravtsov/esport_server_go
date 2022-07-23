package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gym struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	Address   string             `json:"address" bson:"address"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`
}
