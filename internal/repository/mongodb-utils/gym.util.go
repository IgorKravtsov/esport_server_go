package mongodb_utils

import (
	"github.com/IgorKravtsov/esport_server_go/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoGym struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	Address   string             `json:"address" bson:"address"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`
}

func GymToDomain(mongoGym MongoGym) domain.Gym {
	return domain.Gym{
		ID:        mongoGym.ID.Hex(),
		Title:     mongoGym.Title,
		Address:   mongoGym.Address,
		CreatedAt: mongoGym.CreatedAt.Time(),
		CreatedBy: mongoGym.CreatedBy.Hex(),
	}
}

func GymToRepo(g domain.Gym) (*MongoGym, error) {
	var ID primitive.ObjectID
	if g.ID != "" {
		id, err := primitive.ObjectIDFromHex(g.ID)
		if err != nil {
			return nil, err
		}
		ID = id
	}
	createdAt := primitive.NewDateTimeFromTime(g.CreatedAt)
	createdBy, err := primitive.ObjectIDFromHex(g.CreatedBy)
	if err != nil {
		return nil, err
	}

	return &MongoGym{
		ID:        ID,
		Title:     g.Title,
		Address:   g.Address,
		CreatedAt: createdAt,
		CreatedBy: createdBy,
	}, nil
}
