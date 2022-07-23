package domain

import (
	"time"
)

type TrainingType string

const (
	Technique TrainingType = "technique"
	Kata                   = "kata"
)

type Training struct {
	ID               string       `json:"id" bson:"_id,omitempty"`
	TrainingDateTime time.Time    `json:"trainingDateTime" bson:"trainingDateTime"`
	TrainingType     TrainingType `json:"training_type" bson:"training_type"`
	GymID            string       `json:"gym_id" bson:"gym_id"`
}
