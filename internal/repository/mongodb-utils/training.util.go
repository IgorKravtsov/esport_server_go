package mongodb_utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/IgorKravtsov/esport_server_go/internal/domain"
)

type MongoTraining struct {
	ID               primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	TrainingDateTime primitive.DateTime  `json:"trainingDateTime" bson:"trainingDateTime"`
	TrainingType     domain.TrainingType `json:"training_type" bson:"training_type"`
	GymID            primitive.ObjectID  `json:"gym_id" bson:"gym_id"`
}

func TrainingToDomain(mongoTraining MongoTraining) domain.Training {
	return domain.Training{
		ID:               mongoTraining.ID.Hex(),
		TrainingDateTime: mongoTraining.TrainingDateTime.Time(),
		TrainingType:     mongoTraining.TrainingType,
		GymID:            mongoTraining.GymID.Hex(),
	}
}

func TrainingToRepo(t domain.Training) (*MongoTraining, error) {
	var ID primitive.ObjectID
	if t.ID != "" {
		id, err := primitive.ObjectIDFromHex(t.ID)
		if err != nil {
			return nil, err
		}
		ID = id
	}
	trainingTime := primitive.NewDateTimeFromTime(t.TrainingDateTime)
	gymID, err := primitive.ObjectIDFromHex(t.GymID)
	if err != nil {
		return nil, err
	}

	return &MongoTraining{
		ID:               ID,
		TrainingDateTime: trainingTime,
		TrainingType:     t.TrainingType,
		GymID:            gymID,
	}, nil
}
