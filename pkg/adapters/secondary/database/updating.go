package database

import (
	"database/sql"
	"errors"
)

func (a *dbadapter) CloseDBConn() error {
	return a.db.Close()
}

func (a *dbadapter) UpdateEmail(uuid string, email string) error {
	var (
		tx      *sql.Tx
		result  sql.Result
		numRows int64
		err     error
	)

	tx, err = a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err = a.statements[UPDATE_EMAIL].Exec(email, uuid)
	if err != nil {
		return err
	}

	numRows, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if numRows == 0 {
		return errors.New("notfound")
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
