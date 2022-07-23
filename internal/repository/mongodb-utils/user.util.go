package mongodb_utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/IgorKravtsov/esport_server_go/internal/domain"
)

type MongoUser struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name         string               `json:"name" bson:"name"`
	Email        string               `json:"email" bson:"email"`
	Phone        string               `json:"phone" bson:"phone"`
	Password     string               `json:"password" bson:"password"`
	RegisteredAt primitive.DateTime   `json:"registeredAt" bson:"registeredAt"`
	LastVisitAt  primitive.DateTime   `json:"lastVisitAt" bson:"lastVisitAt"`
	Verification domain.Verification  `json:"verification" bson:"verification"`
	Trainers     []primitive.ObjectID `json:"trainers" bson:"trainers"`
}

func UserToDomain(mongoUser MongoUser) domain.User {
	var trainerIds []string
	for _, objId := range mongoUser.Trainers {
		trainerIds = append(trainerIds, objId.Hex())
	}
	return domain.User{
		ID:           mongoUser.ID.Hex(),
		Name:         mongoUser.Name,
		Email:        mongoUser.Email,
		Phone:        mongoUser.Phone,
		Password:     mongoUser.Password,
		RegisteredAt: mongoUser.RegisteredAt.Time(),
		LastVisitAt:  mongoUser.LastVisitAt.Time(),
		Verification: mongoUser.Verification,
		Trainers:     trainerIds,
	}
}

func UserToRepo(u domain.User) (*MongoUser, error) {
	var ID primitive.ObjectID
	if u.ID != "" {
		id, err := primitive.ObjectIDFromHex(u.ID)
		if err != nil {
			return nil, err
		}
		ID = id
	}

	var trainerIds []primitive.ObjectID
	for _, id := range u.Trainers {
		trainerID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		trainerIds = append(trainerIds, trainerID)
	}

	registeredAt := primitive.NewDateTimeFromTime(u.RegisteredAt)
	lastVisitAt := primitive.NewDateTimeFromTime(u.LastVisitAt)

	return &MongoUser{
		ID:           ID,
		Name:         u.Name,
		Email:        u.Email,
		Phone:        u.Phone,
		Password:     u.Password,
		RegisteredAt: registeredAt,
		LastVisitAt:  lastVisitAt,
		Verification: u.Verification,
		Trainers:     trainerIds,
	}, nil
}
