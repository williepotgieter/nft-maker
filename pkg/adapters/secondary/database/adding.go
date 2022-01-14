package database

import (
	"database/sql"

	"github.com/williepotgieter/nft-maker/pkg/domain/models"
)

func (a *dbadapter) CreateUser(user models.User) error {
	var (
		tx  *sql.Tx
		err error
	)

	tx, err = a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	defer a.statements[CREATE_USER].Close()

	_, err = a.statements[CREATE_USER].Exec(
		user.CreatedAt,
		user.ModifiedAt,
		user.UUID,
		user.Name,
		user.Surname,
		user.Email,
		user.Password,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
