package repository

import (
	"fmt"

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
	db.AutoMigrate(&Person{}, &Course{}, &Favorite{}, &Survey{})
}

// CloseDB ...
func CloseDB() {
	defer db.Close()
}
