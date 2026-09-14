package config

import (
	"context"
	amqpConfig "github.com/mephistolie/chefbook-backend-common/mq/config"
	"github.com/mephistolie/chefbook-backend-subscription/internal/logging"
)

const (
	EnvDev  = "develop"
	EnvProd = "production"
)

type Config struct {
	Environment *string
	Port        *int
	LogsPath    *string

	Google Google

	Firebase Firebase

	AuthService AuthService

	Database Database
	Amqp     amqpConfig.Amqp
	Smtp     Smtp
}

type Google struct {
	PackageName      *string
	JsonKey          *string
	ProductIdPremium *string
	ProductIdMaximum *string
}

type Firebase struct {
	Credentials *string
}

type AuthService struct {
	Addr *string
}

type Database struct {
	Host     *string
	Port     *int
	User     *string
	Password *string
	DBName   *string
}

type Smtp struct {
	Username     *string
	Host         *string
	Port         *int
	Email        *string
	Password     *string
	SendAttempts *int
}

func (c Config) Validate() error {
	if *c.Environment != EnvProd {
		*c.Environment = EnvDev
	}
	return nil
}

func (c Config) Print(ctx context.Context) {
	logging.Events{}.ConfigLoaded(ctx)
}
