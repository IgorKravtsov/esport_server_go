package gym

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/IgorKravtsov/esport_server_go/internal/domain"
	"github.com/IgorKravtsov/esport_server_go/internal/repository"
	"github.com/IgorKravtsov/esport_server_go/internal/service/gym/dto"
)

type Gym interface {
	Create(ctx context.Context, input dto.CreateGym, creatorID primitive.ObjectID) error
}

type Service struct {
	repo   repository.Gym
	domain string
}

func NewGymService(
	repo repository.Gym, domain string) *Service {
	return &Service{
		repo:   repo,
		domain: domain,
	}
}

func (s Service) Create(ctx context.Context, input dto.CreateGym, creatorID primitive.ObjectID) error {
	gym := domain.Gym{
		Title:   input.Title,
		Address: input.Address,
	}
	if err := s.repo.Create(ctx, gym); err != nil {
		return err
	}

	return nil
}