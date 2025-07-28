package main

import (
	"go-mini-server/core/db/pool"
	"go-mini-server/internal/storage"
	"go-mini-server/internal/user"
)

var (
	userService user.Service
)

func initServices(c Config) {
	mainPool := pool.NewService(c.DB)

	userRepository := storage.NewUserRepository(mainPool)
	userService = user.NewService(userRepository)
}
