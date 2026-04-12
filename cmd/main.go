package main

import (
	"practice-7/internal/controller/http/v1"
	"practice-7/internal/usecase"
	"practice-7/internal/usecase/repo"
	"practice-7/pkg/postgres"
	"practice-7/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	pg, _ := postgres.New()

	userRepo := repo.NewUserRepo(pg)
	userUseCase := usecase.NewUserUseCase(userRepo)

	r := gin.Default()

	r.Use(utils.RateLimiter())

	api := r.Group("/api")
	v1.NewUserRoutes(api, userUseCase)

	r.Run(":8080")
}
