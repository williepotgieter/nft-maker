package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBConn struct {
	db *gorm.DB
}

func NewDBConn(file string) *DBConn {

	var (
		db  *gorm.DB
		err error
	)

	log.Println("Connecting to database...")
	db, err = gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	log.Println("Successfully connected to database!")

	// Migrate db schema
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to automigrate database")
	}
	log.Println("Database automigration completed...")

	return &DBConn{db}
}
