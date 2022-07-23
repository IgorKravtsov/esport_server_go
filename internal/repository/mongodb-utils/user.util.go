package mongodb_utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/IgorKravtsov/esport_server_go/internal/domain"
)

type MongoUser struct {
	ID           primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name         string              `json:"name" bson:"name"`
	Email        string              `json:"email" bson:"email"`
	Phone        string              `json:"phone" bson:"phone"`
	Password     string              `json:"password" bson:"password"`
	RegisteredAt primitive.DateTime  `json:"registeredAt" bson:"registeredAt"`
	LastVisitAt  primitive.DateTime  `json:"lastVisitAt" bson:"lastVisitAt"`
	Verification domain.Verification `json:"verification" bson:"verification"`
	//Schools []primitive.ObjectID `json:"schools" bson:"schools"`
}

func UserToDomain(mongoUser MongoUser) domain.User {
	return domain.User{
		ID:           mongoUser.ID.Hex(),
		Name:         mongoUser.Name,
		Email:        mongoUser.Email,
		Phone:        mongoUser.Phone,
		Password:     mongoUser.Password,
		RegisteredAt: mongoUser.RegisteredAt.Time(),
		LastVisitAt:  mongoUser.LastVisitAt.Time(),
		Verification: mongoUser.Verification,
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
	}, nil
}
