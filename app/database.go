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
	err = db.AutoMigrate(&User{}, &Account{})
	if err != nil {
		panic("failed to automigrate database")
	}

	if res := db.Exec("PRAGMA foreign_keys = ON", nil); res.Error != nil {
		panic(res.Error)
	}

	log.Println("Database automigration completed...")

	return &DBConn{db}
}
