package main

import (
	"flag"
	amqpConfig "github.com/mephistolie/chefbook-backend-common/mq/config"
	"github.com/mephistolie/chefbook-backend-subscription/internal/app"
	"github.com/mephistolie/chefbook-backend-subscription/internal/config"
	"github.com/peterbourgon/ff/v3"
	"os"
)

func main() {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	cfg := config.Config{
		Environment: fs.String("environment", "debug", "service environment"),
		Port:        fs.Int("port", 8080, "service port"),
		LogsPath:    fs.String("logs-path", "", "logs file path"),

		Google: config.Google{
			PackageName:      fs.String("google-package-name", "", "Google App package name"),
			JsonKey:          fs.String("google-json-key", "", "Google JSON Key for validating Google subscriptions"),
			ProductIdPremium: fs.String("google-product-id-premium", "premium", "Google Premium Product ID"),
			ProductIdMaximum: fs.String("google-product-id-maximum", "maximum", "Google Maximum Product ID"),
		},

		Firebase: config.Firebase{
			Credentials: fs.String("firebase-credentials", "", "Firebase credentials JSON; leave empty to disable"),
		},

		AuthService: config.AuthService{
			Addr: fs.String("auth-addr", "", "auth service address"),
		},

		Database: config.Database{
			Host:     fs.String("db-host", "localhost", "database host"),
			Port:     fs.Int("db-port", 5432, "database port"),
			User:     fs.String("db-user", "", "database user name"),
			Password: fs.String("db-password", "", "database user password"),
			DBName:   fs.String("db-name", "", "service database name"),
		},

		Amqp: amqpConfig.Amqp{
			Host:     fs.String("amqp-host", "", "message broker host; leave empty to disable"),
			Port:     fs.Int("amqp-port", 5672, "message broker port"),
			User:     fs.String("amqp-user", "guest", "message broker user name"),
			Password: fs.String("amqp-password", "guest", "message broker user password"),
			VHost:    fs.String("amqp-vhost", "", "message broker virtual host"),
		},

		Smtp: config.Smtp{
			Host:         fs.String("smtp-host", "", "SMTP host; leave empty to disable"),
			Port:         fs.Int("smtp-port", 465, "SMTP port"),
			Email:        fs.String("smtp-email", "", "SMTP sender email"),
			Password:     fs.String("smtp-password", "", "SMTP sender password"),
			SendAttempts: fs.Int("smtp-attempts", 3, "SMTP email sending attempts"),
		},
	}
	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVars()); err != nil {
		panic(err)
	}

	err := cfg.Validate()
	if err != nil {
		panic(err)
	}

	app.Run(&cfg)
}
