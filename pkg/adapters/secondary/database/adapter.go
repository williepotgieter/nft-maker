package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// SQL Statement names
	CREATE_USER   = "create_user"
	GET_ALL_USERS = "get_all_users"
	GET_USER      = "get_user"
)

type dbadapter struct {
	db         *sql.DB
	statements map[string]*sql.Stmt
}

func NewDBAdapter() *dbadapter {
	var (
		db    *sql.DB
		stmts = make(map[string]*sql.Stmt)
		uri   string
		err   error
	)

	uri = fmt.Sprintf("%s:%s@tcp(nft-maker-db:3306)/%s", os.Getenv("MARIADB_USER"), os.Getenv("MARIADB_PASSWORD"), os.Getenv("MARIADB_NAME"))

	db, err = sql.Open("mysql", uri)
	if err != nil {
		log.Fatalln("unable to connect to database")
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln("db connection failed", err)
	} else {
		log.Println("connected to database...")
	}

	stmts = PrepareStatements(db)

	return &dbadapter{db, stmts}
}

func PrepareStatements(db *sql.DB) map[string]*sql.Stmt {
	var (
		prepStmts = make(map[string]*sql.Stmt)
		err       error
	)

	// Create a new user
	prepStmts[CREATE_USER], err = db.Prepare(`
	INSERT INTO users (created_at, modified_at, uuid, name, surname, email, password)
	VALUES (?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		log.Fatalln(err)
	}

	// Get all users
	prepStmts[GET_ALL_USERS], err = db.Prepare(`
	SELECT uuid, name, surname, email, created_at, modified_at FROM users;`)
	if err != nil {
		log.Fatalln(err)
	}

	// Get specific user
	prepStmts[GET_USER], err = db.Prepare(`
	SELECT uuid, name, surname, email, created_at, modified_at FROM users WHERE uuid = ?;`)
	if err != nil {
		log.Fatalln(err)
	}

	return prepStmts
}
