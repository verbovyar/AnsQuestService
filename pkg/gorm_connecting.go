package pkg

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(connection_string string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(connection_string), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return db
}
