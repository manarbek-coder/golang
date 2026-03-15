package main

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"Assignment5/internal/handlers"
	"Assignment5/internal/repository/users"
	"Assignment5/internal/usecase"
)

func main() {

	db, err := sqlx.Connect("postgres", "postgres://postgres:postgresql@localhost:5432/mydb?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	repo := users.NewRepository(db)
	usecase := usecase.NewUserUsecase(repo)
	handler := handlers.NewUserHandler(usecase)

	http.HandleFunc("/health", handler.Health)
	http.HandleFunc("/users", handler.GetUsers)
	http.HandleFunc("/common-friends", handler.CommonFriends)

	log.Println("server running :8080")

	http.ListenAndServe(":8080", nil)
}
