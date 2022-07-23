package repository

import (
	"context"
	"errors"
	"github.com/IgorKravtsov/esport_server_go/internal/domain"
	"github.com/IgorKravtsov/esport_server_go/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const usersCollection = "users"

// User interface
type User interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	Verify(ctx context.Context, userID primitive.ObjectID, code string) error
	//SetSession(ctx context.Context, userID primitive.ObjectID, session domain.Session) error
	//AttachSchool(ctx context.Context, userID, schoolID primitive.ObjectID) error
}

type UserRepo struct {
	db *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{
		db: db.Collection(usersCollection),
	}
}

func (r *UserRepo) Create(ctx context.Context, user domain.User) error {
	_, err := r.db.InsertOne(ctx, user)
	if mongodb.IsDuplicate(err) {
		return domain.ErrUserAlreadyExists
	}

	return err
}

func (r *UserRepo) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User
	if err := r.db.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	var user domain.User
	if err := r.db.FindOne(ctx, bson.M{
		"session.refreshToken": refreshToken,
		"session.expiresAt":    bson.M{"$gt": time.Now()},
	}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepo) Verify(ctx context.Context, userID primitive.ObjectID, code string) error {
	res, err := r.db.UpdateOne(ctx,
		bson.M{"verification.code": code, "_id": userID},
		bson.M{"$set": bson.M{"verification.verified": true, "verification.code": ""}})
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return domain.ErrVerificationCodeInvalid
	}

	return nil
}

//func (r *UserRepo) SetSession(ctx context.Context, userID primitive.ObjectID, session domain.Session) error {
//  _, err := r.db.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}})
//
//  return err
//}
//
//func (r *UserRepo) AttachSchool(ctx context.Context, userID, schoolId primitive.ObjectID) error {
//  _, err := r.db.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$push": bson.M{"schools": schoolId}})
//
//  return err
//}
