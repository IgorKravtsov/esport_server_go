package app

import (
  "backend/pkg/auth"
  "backend/pkg/database/mongodb"
  "backend/pkg/logger"
  "context"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "time"
)

func Run(configPath string) {
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
  
  memCache := cache.NewMemoryCache()
  hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)
  
  emailSender, err := smtp.NewSMTPSender(cfg.SMTP.From, cfg.SMTP.Pass, cfg.SMTP.Host, cfg.SMTP.Port)
  if err != nil {
    logger.Error(err)
    
    return
  }
  
  tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
  if err != nil {
    logger.Error(err)
    
    return
  }
  
  otpGenerator := otp.NewGOTPGenerator()
  
  storageProvider, err := newStorageProvider(cfg)
  if err != nil {
    logger.Error(err)
    
    return
  }
  
  cloudflareClient, err := cloudflare.New(cfg.Cloudflare.ApiKey, cfg.Cloudflare.Email)
  if err != nil {
    logger.Error(err)
    
    return
  }
  
  dnsService := dns.NewService(cloudflareClient, cfg.Cloudflare.ZoneEmail, cfg.Cloudflare.CnameTarget)
  
  // Services, Repos & API Handlers
  repos := repository.NewRepositories(db)
  services := service.NewServices(service.Deps{
    Repos:                  repos,
    Cache:                  memCache,
    Hasher:                 hasher,
    TokenManager:           tokenManager,
    EmailSender:            emailSender,
    EmailConfig:            cfg.Email,
    AccessTokenTTL:         cfg.Auth.JWT.AccessTokenTTL,
    RefreshTokenTTL:        cfg.Auth.JWT.RefreshTokenTTL,
    FondyCallbackURL:       cfg.Payment.FondyCallbackURL,
    CacheTTL:               int64(cfg.CacheTTL.Seconds()),
    OtpGenerator:           otpGenerator,
    VerificationCodeLength: cfg.Auth.VerificationCodeLength,
    StorageProvider:        storageProvider,
    Environment:            cfg.Environment,
    Domain:                 cfg.HTTP.Host,
    DNS:                    dnsService,
  })
  handlers := delivery.NewHandler(services, tokenManager)
  
  services.Files.InitStorageUploaderWorkers(context.Background())
  
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
