package main

import (
	"github.com/dewadg/ecdh-exchange/internal/app/server"
	"github.com/dewadg/ecdh-exchange/internal/pkg/config"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("error loading config")
	}

	server.Run(cfg)
}
