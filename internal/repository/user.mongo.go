package repository

import (
	"context"
	"errors"
	"github.com/IgorKravtsov/esport_server_go/internal/domain"
	mongodb_utils "github.com/IgorKravtsov/esport_server_go/internal/repository/mongodb-utils"
	"github.com/IgorKravtsov/esport_server_go/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const usersCollection = "users"

type Verification struct {
	Code     string `json:"code" bson:"code"`
	Verified bool   `json:"verified" bson:"verified"`
}

// User interface
type User interface {
	Create(ctx context.Context, user domain.User) (string, error)
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	Verify(ctx context.Context, userID string, code string) error
	SetSession(ctx context.Context, userID string, session domain.Session) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
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

func (r *UserRepo) Create(ctx context.Context, u domain.User) (string, error) {
	user, err := mongodb_utils.UserToRepo(u)
	if err != nil {
		return "", err
	}
	insertResult, err := r.db.InsertOne(ctx, user)
	if mongodb.IsDuplicate(err) {
		return "", domain.ErrUserAlreadyExists
	}
	oid, err := mongodb_utils.ConvertToObjID(insertResult.InsertedID)
	if err != nil {
		return "", err
	}

	return oid.Hex(), nil
}

func (r *UserRepo) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var user mongodb_utils.MongoUser
	if err := r.db.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	return mongodb_utils.UserToDomain(user), nil
}

func (r *UserRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	var user mongodb_utils.MongoUser
	if err := r.db.FindOne(ctx, bson.M{
		"session.refreshToken": refreshToken,
		"session.expiresAt":    bson.M{"$gt": time.Now()},
	}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	return mongodb_utils.UserToDomain(user), nil
}

func (r *UserRepo) Verify(ctx context.Context, userID string, code string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err == nil {
		return err
	}
	res, err := r.db.UpdateOne(ctx,
		bson.M{"verification.code": code, "_id": objID},
		bson.M{"$set": bson.M{"verification.verified": true, "verification.code": ""}})
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return domain.ErrVerificationCodeInvalid
	}

	return nil
}

func (r *UserRepo) SetSession(ctx context.Context, userID string, session domain.Session) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err == nil {
		return err
	}
	_, err = r.db.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}})

	return err
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user mongodb_utils.MongoUser
	if err := r.db.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}
	domainUser := mongodb_utils.UserToDomain(user)
	return &domainUser, nil
}

//
//func (r *UserRepo) AttachSchool(ctx context.Context, userID, schoolId primitive.ObjectID) error {
//  _, err := r.db.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$push": bson.M{"schools": schoolId}})
//
//  return err
//}
