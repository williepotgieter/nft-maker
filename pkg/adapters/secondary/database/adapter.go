package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type dbadapter struct {
	db *sql.DB
}

func NewDBAdapter() *dbadapter {
	var (
		db  *sql.DB
		uri string
		err error
	)

	uri = fmt.Sprintf("root:%s@tcp(db:3306)/nftmaker", os.Getenv("MYSQL_ROOT_PASSWORD"))

	db, err = sql.Open("mysql", uri)
	if err != nil {
		log.Fatalln("unable to connect to database")
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln("db connection failed")
	} else {
		log.Println("connected to database...")
	}

	return &dbadapter{db}
}
