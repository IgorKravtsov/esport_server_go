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
	"github.com/IgorKravtsov/esport_server_go/pkg/otp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User interface {
	Register(ctx context.Context, input dto.UserRegister) error
	Login(ctx context.Context, input dto.UserLogin) (tokens.Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (tokens.Tokens, error)
	Verify(ctx context.Context, userID primitive.ObjectID, hash string) error
}

type Service struct {
	repo                   repository.User
	hasher                 hash.PasswordHasher
	tokenManager           auth.TokenManager
	otpGenerator           otp.Generator
	accessTokenTTL         time.Duration
	refreshTokenTTL        time.Duration
	verificationCodeLength int

	domain string
	//dnsService   dns.DomainManager
	//
	//emailService  Emails
	//schoolService Schools

}

func NewUserService(
	repo repository.User, hasher hash.PasswordHasher,
	tokenManager auth.TokenManager,
	accessTTL, refreshTTL time.Duration,
	verificationCodeLength int, domain string, otpGenerator otp.Generator) *Service {
	return &Service{
		repo:                   repo,
		hasher:                 hasher,
		tokenManager:           tokenManager,
		accessTokenTTL:         accessTTL,
		refreshTokenTTL:        refreshTTL,
		verificationCodeLength: verificationCodeLength,
		domain:                 domain,
		otpGenerator:           otpGenerator,
		//emailService:           emailService,
		//schoolService:          schoolsService,
		//dnsService:             dnsService,
	}
}

func (s *Service) Register(ctx context.Context, input dto.UserRegister) error {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	verificationCode := s.otpGenerator.RandomSecret(s.verificationCodeLength)

	user := domain.User{
		Name:         input.Name,
		Password:     passwordHash,
		Email:        input.Email,
		RegisteredAt: time.Now(),
		LastVisitAt:  time.Now(),
		Verification: domain.Verification{
			Code: verificationCode,
		},
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

func (s *Service) Login(ctx context.Context, input dto.UserLogin) (tokens.Tokens, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return tokens.Tokens{}, err
	}

	user, err := s.repo.GetByCredentials(ctx, input.Email, passwordHash)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return tokens.Tokens{}, err
		}

		return tokens.Tokens{}, err
	}

	return s.createSession(ctx, user.ID)
}

func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (tokens.Tokens, error) {
	user, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return tokens.Tokens{}, err
	}

	return s.createSession(ctx, user.ID)
}

func (s *Service) Verify(ctx context.Context, userID primitive.ObjectID, hash string) error {
	err := s.repo.Verify(ctx, userID, hash)
	if err != nil {
		if errors.Is(err, domain.ErrVerificationCodeInvalid) {
			return err
		}

		return err
	}

	return nil
}

func (s *Service) createSession(ctx context.Context, userId primitive.ObjectID) (tokens.Tokens, error) {
	var (
		res tokens.Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userId.Hex(), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.repo.SetSession(ctx, userId, session)

	return res, err
}
