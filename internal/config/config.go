package config

import (
	"errors"
	"log"
	"os"

	"github.com/sakamoto-max/ratelimiter/internal/env"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	STAGE    string
	Grpc     GrpcConfig
	Redis    RedisConfig
	Postgres PostgresConfig
	Http     HTTPConfig
	Auth     AuthConfig
}

type GrpcConfig struct {
	Port string `validate:"required"`
}

type HTTPConfig struct {
	Port string `validate:"required"`
}

type RedisConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	Db       string `validate:"required"`
	UserName string `validate:"required"`
	Password string `validate:"required"`
}

type PostgresConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	Db       string `validate:"required"`
	UserName string `validate:"required"`
	Password string `validate:"required"`
	SSLmode  string `validate:"required"`
}

type AuthConfig struct {
	SecretKey string `validate:"required"`
}

func New() *Config {

	env.LoadEnv("../../app.env")

	stage := os.Getenv("STAGE")
	if stage == "" {
		log.Fatalf("failed to get : %v from env", "STAGE")
	}

	grpcConfig := GrpcConfig{
		Port: os.Getenv("GRPC_PORT"),
	}

	redisConfig := RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Db:       os.Getenv("REDIS_DB"),
		UserName: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	postgresConfig := PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Db:       os.Getenv("POSTGRES_DB"),
		UserName: os.Getenv("POSTGRES_USERNAME"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		SSLmode:  os.Getenv("POSTGRES_SSL_MODE"),
	}

	httpConfig := HTTPConfig{
		Port: os.Getenv("HTTP_PORT"),
	}

	authConfig := AuthConfig{
		SecretKey: os.Getenv("SECRET_KEY"),
	}

	config := Config{
		STAGE:    stage,
		Grpc:     grpcConfig,
		Redis:    redisConfig,
		Postgres: postgresConfig,
		Http:     httpConfig,
		Auth:     authConfig,
	}

	newValidator := validator.New()

	err := newValidator.Struct(config)
	if err != nil {
		var validatorErr validator.ValidationErrors
		if errors.As(err, &validatorErr) {
			log.Fatalf("failed to get : %v from env", validatorErr.Error())
		}
	}

	return &config
}
