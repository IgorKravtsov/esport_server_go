package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type TrainingType string

const (
	Technique TrainingType = "technique"
	Kata                   = "kata"
)

type Training struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TrainingDateTime primitive.DateTime `json:"trainingDateTime" bson:"trainingDateTime"`
	TrainingType     TrainingType       `json:"training_type" bson:"training_type"`
	GymID            primitive.ObjectID `json:"gym_id" bson:"gym_id"`
}
