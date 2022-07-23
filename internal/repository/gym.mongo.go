package repository

import (
	"context"
	"github.com/IgorKravtsov/esport_server_go/internal/domain"
	"github.com/IgorKravtsov/esport_server_go/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

const gymsCollection = "gyms"

// Gym interface
type Gym interface {
	Create(ctx context.Context, user domain.Gym) error
	//GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	//GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	//Verify(ctx context.Context, userID primitive.ObjectID, code string) error
	//SetSession(ctx context.Context, userID primitive.ObjectID, session domain.Session) error
	//AttachSchool(ctx context.Context, userID, schoolID primitive.ObjectID) error
}

type GymRepo struct {
	db *mongo.Collection
}

func NewGymRepo(db *mongo.Database) *GymRepo {
	return &GymRepo{
		db: db.Collection(gymsCollection),
	}
}

func (r *GymRepo) Create(ctx context.Context, gym domain.Gym) error {
	_, err := r.db.InsertOne(ctx, gym)
	if mongodb.IsDuplicate(err) {
		return domain.ErrorAlreadyExitsts(gymsCollection)
	}

	return err
}
