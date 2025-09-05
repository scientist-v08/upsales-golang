package initializers

import (
	"log"
	"os"

	"github.com/scientist-v08/favmovies/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {
    var err error

    dsn := os.Getenv("DB_URL")
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        log.Fatal("Failed to connect to database")
    }

	// AutoMigrate creates the table if it doesn't exist
	err = DB.AutoMigrate(&model.Movies{}, &model.User{})
	if err != nil {
    	log.Fatal("Failed to migrate database: ", err)
	}
}