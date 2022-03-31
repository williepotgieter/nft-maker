package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBConn struct {
	db *gorm.DB
}

func NewDBConn(file string) *DBConn {
	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate db schema
	db.AutoMigrate(&User{})

	return &DBConn{db}
}
