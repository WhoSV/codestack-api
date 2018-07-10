package database

import (
	"fmt"

	"github.com/WhoSV/codestack-api/model"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// ConnectDB ...
func ConnectDB() {
	var err error
	db, err = gorm.Open("sqlite3", "codestack.db")
	if err != nil {
		fmt.Printf("database connection failed: %v", err)
		return
	}
	db.AutoMigrate(&model.Person{}, &model.Course{}, &model.Favorite{}, &model.Survey{})
	// defer db.Close()
}

// DB returns the current db instance.
func DB() *gorm.DB {
	return db
}
