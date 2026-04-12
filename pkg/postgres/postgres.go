package postgres

import (
	"practice-7/internal/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func New() (*Postgres, error) {
	dsn := "host=localhost user=postgres password=postgresql dbname=practice7 port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entity.User{})

	return &Postgres{DB: db}, nil
}
