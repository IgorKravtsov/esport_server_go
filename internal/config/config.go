package config

import (
  "github.com/spf13/viper"
  "os"
  "time"
)

const (
  defaultHTTPPort               = "8000"
  defaultHTTPRWTimeout          = 10 * time.Second
  defaultHTTPMaxHeaderMegabytes = 1
  defaultAccessTokenTTL         = 15 * time.Minute
  defaultRefreshTokenTTL        = 24 * time.Hour * 30
  //defaultLimiterRPS             = 10
  //defaultLimiterBurst           = 2
  //defaultLimiterTTL             = 10 * time.Minute
  defaultVerificationCodeLength = 8
  
  EnvLocal = "local"
  Prod     = "prod"
)

type (
  Config struct {
    Environment string
    Mongo       MongoConfig
    HTTP        HTTPConfig
    Auth        AuthConfig
  }
  
  MongoConfig struct {
    URI      string
    User     string
    Password string
    Name     string `mapstructure:"databaseName"`
  }
  
  HTTPConfig struct {
    Host               string        `mapstructure:"host"`
    Port               string        `mapstructure:"port"`
    ReadTimeout        time.Duration `mapstructure:"readTimeout"`
    WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
    MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
  }
  
  AuthConfig struct {
    JWT                    JWTConfig
    PasswordSalt           string
    VerificationCodeLength int `mapstructure:"verificationCodeLength"`
  }
  
  JWTConfig struct {
    AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
    RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
    SigningKey      string
  }
)

func Init(configsDir string) (*Config, error) {
  populateDefaults()
  
  if err := parseConfigFile(configsDir, os.Getenv("APP_ENV")); err != nil {
    return nil, err
  }
  
  var cfg Config
  if err := unmarshal(&cfg); err != nil {
    return nil, err
  }
  
  setFromEnv(&cfg)
  
  return &cfg, nil
}

func parseConfigFile(folder, env string) error {
  viper.AddConfigPath(folder)
  viper.SetConfigName("main")
  
  if err := viper.ReadInConfig(); err != nil {
    return err
  }
  
  if env == EnvLocal {
    return nil
  }
  
  viper.SetConfigName(env)
  
  return viper.MergeInConfig()
}

func populateDefaults() {
  viper.SetDefault("http.port", defaultHTTPPort)
  viper.SetDefault("http.max_header_megabytes", defaultHTTPMaxHeaderMegabytes)
  viper.SetDefault("http.timeouts.read", defaultHTTPRWTimeout)
  viper.SetDefault("http.timeouts.write", defaultHTTPRWTimeout)
  viper.SetDefault("auth.accessTokenTTL", defaultAccessTokenTTL)
  viper.SetDefault("auth.refreshTokenTTL", defaultRefreshTokenTTL)
  viper.SetDefault("auth.verificationCodeLength", defaultVerificationCodeLength)
  //viper.SetDefault("limiter.rps", defaultLimiterRPS)
  //viper.SetDefault("limiter.burst", defaultLimiterBurst)
  //viper.SetDefault("limiter.ttl", defaultLimiterTTL)
}
