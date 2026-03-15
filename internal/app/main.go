package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"Assignment3/internal/repository/postgres"
	"Assignment3/internal/repository/users"
	"Assignment3/internal/usecase"
	"Assignment3/pkg/modules"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := initPostgreSQL()

	postgres.AutoMigrate(dbConfig)

	pgDialect := postgres.NewPGXDialect(ctx, dbConfig)

	userRepo := users.NewRepository(pgDialect.DB)

	userUsecase := usecase.NewUserUsecase(userRepo)

	allUsers, err := userUsecase.GetAllUsers()
	if err != nil {
		log.Printf("Error fetching all users: %v", err)
	} else {
		fmt.Printf("All users: %+v\n", allUsers)
	}

	user, err := userUsecase.GetUserByID(1)
	if err != nil {
		log.Printf("Error fetching user by ID: %v", err)
	} else {
		fmt.Printf("User by ID 1: %+v\n", user)
	}

	fmt.Println("Database layer initialized successfully")
}

func initPostgreSQL() *modules.PostgreSQLConfig {
	return &modules.PostgreSQLConfig{
		Host:        "localhost",
		Port:        "5432",
		Username:    "postgres",
		Password:    "postgresql",
		DBName:      "mydb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}
}
