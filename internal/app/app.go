package app

import (
	"context"
	"errors"
	"github.com/IgorKravtsov/esport_server_go/internal/server"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IgorKravtsov/esport_server_go/internal/config"
	delivery "github.com/IgorKravtsov/esport_server_go/internal/delivery/http"
	"github.com/IgorKravtsov/esport_server_go/internal/repository"
	"github.com/IgorKravtsov/esport_server_go/internal/service"
	"github.com/IgorKravtsov/esport_server_go/pkg/auth"
	"github.com/IgorKravtsov/esport_server_go/pkg/database/mongodb"
	"github.com/IgorKravtsov/esport_server_go/pkg/hash"
	"github.com/IgorKravtsov/esport_server_go/pkg/logger"
)

// @title Creatly API
// @version 1.0
// @description REST API for Creatly App

// @host localhost:8000
// @BasePath /api/v1/

// @securityDefinitions.apikey AdminAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey StudentsAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey UsersAuth
// @in header
// @name Authorization

// Run initializes whole application.
func Run(configPath string) {
	if err := godotenv.Load(); err != nil {
		logger.Errorf("error loading env variables: %s", err.Error())
	}

	cfg, err := config.Init(configPath)
	if err != nil {
		logger.Error(err)

		return
	}

	// Dependencies
	mongoClient, err := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)
	if err != nil {
		logger.Error(err)

		return
	}

	db := mongoClient.Database(cfg.Mongo.Name)

	//memCache := cache.NewMemoryCache()
	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)

	//emailSender, err := smtp.NewSMTPSender(cfg.SMTP.From, cfg.SMTP.Pass, cfg.SMTP.Host, cfg.SMTP.Port)
	//if err != nil {
	//  logger.Error(err)
	//
	//  return
	//}

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)

		return
	}

	//otpGenerator := otp.NewGOTPGenerator()
	//
	//storageProvider, err := newStorageProvider(cfg)
	//if err != nil {
	//  logger.Error(err)
	//
	//  return
	//}

	//cloudflareClient, err := cloudflare.New(cfg.Cloudflare.ApiKey, cfg.Cloudflare.Email)
	//if err != nil {
	//  logger.Error(err)
	//
	//  return
	//}

	//dnsService := dns.NewService(cloudflareClient, cfg.Cloudflare.ZoneEmail, cfg.Cloudflare.CnameTarget)

	// Services, Repos & API Handlers
	repos := repository.NewRepositories(db)
	services := service.NewServices(service.Deps{
		Repos:                  repos,
		Hasher:                 hasher,
		TokenManager:           tokenManager,
		AccessTokenTTL:         cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL:        cfg.Auth.JWT.RefreshTokenTTL,
		VerificationCodeLength: cfg.Auth.VerificationCodeLength,
		Environment:            cfg.Environment,
		Domain:                 cfg.HTTP.Host,
		//Cache:                  memCache,
		//EmailSender:            emailSender,
		//EmailConfig:            cfg.Email,
		//FondyCallbackURL:       cfg.Payment.FondyCallbackURL,
		//CacheTTL: int64(cfg.CacheTTL.Seconds()),
		//OtpGenerator:           otpGenerator,
		//StorageProvider:        storageProvider,
		//DNS:                    dnsService,
	})
	handlers := delivery.NewHandler(services, tokenManager)

	//services.Files.InitStorageUploaderWorkers(context.Background())

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}

	if err := mongoClient.Disconnect(context.Background()); err != nil {
		logger.Error(err.Error())
	}
}
