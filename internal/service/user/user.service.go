package user

import (
	"context"
	"errors"
	"github.com/IgorKravtsov/esport_server_go/internal/domain"
	"github.com/IgorKravtsov/esport_server_go/internal/repository"
	"github.com/IgorKravtsov/esport_server_go/internal/service/tokens"
	"github.com/IgorKravtsov/esport_server_go/internal/service/user/dto"
	"github.com/IgorKravtsov/esport_server_go/pkg/auth"
	"github.com/IgorKravtsov/esport_server_go/pkg/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Users interface {
	Register(ctx context.Context, input dto.UserRegister) error
	Login(ctx context.Context, input dto.UserLogin) (tokens.Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (tokens.Tokens, error)
	Verify(ctx context.Context, userID primitive.ObjectID, hash string) error
	//CreateSchool(ctx context.Context, userID primitive.ObjectID, schoolName string) (domain.School, error)
}

type Dto struct {
	Login    dto.UserLogin
	Register dto.UserRegister
}

type UsersService struct {
	repo         repository.Users
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
	//otpGenerator otp.Generator
	//dnsService   dns.DomainManager
	//
	//emailService  Emails
	//schoolService Schools

	accessTokenTTL         time.Duration
	refreshTokenTTL        time.Duration
	verificationCodeLength int

	domain string
}

func NewUsersService(
	repo repository.Users, hasher hash.PasswordHasher,
	tokenManager auth.TokenManager,
	accessTTL, refreshTTL time.Duration,
	verificationCodeLength int, domain string) *UsersService {
	return &UsersService{
		repo:                   repo,
		hasher:                 hasher,
		tokenManager:           tokenManager,
		accessTokenTTL:         accessTTL,
		refreshTokenTTL:        refreshTTL,
		verificationCodeLength: verificationCodeLength,
		domain:                 domain,
		//emailService:           emailService,
		//schoolService:          schoolsService,
		//otpGenerator:           otpGenerator,
		//dnsService:             dnsService,
	}
}

func (s *UsersService) Register(ctx context.Context, input dto.UserRegister) error {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	//verificationCode := s.otpGenerator.RandomSecret(s.verificationCodeLength)

	user := domain.User{
		Name:         input.Name,
		Password:     passwordHash,
		Email:        input.Email,
		RegisteredAt: time.Now(),
		LastVisitAt:  time.Now(),
		//Verification: domain.Verification{
		//  Code: verificationCode,
		//},
	}

	if err = s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return err
		}

		return err
	}

	return nil
	// todo. DECIDE ON EMAIL MARKETING STRATEGY

	//return s.emailService.SendUserVerificationEmail(VerificationEmailInput{
	//  Email:            user.Email,
	//  Name:             user.Name,
	//  VerificationCode: verificationCode,
	//})
}

func (s *UsersService) Login(ctx context.Context, input dto.UserLogin) (tokens.Tokens, error) {
	panic("implement me")
}

func (s *UsersService) RefreshTokens(ctx context.Context, refreshToken string) (tokens.Tokens, error) {
	panic("implement me")
}

func (s *UsersService) Verify(ctx context.Context, userID primitive.ObjectID, hash string) error {
	panic("implement me")
}
