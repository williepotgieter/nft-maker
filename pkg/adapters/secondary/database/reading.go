package database

import (
	"database/sql"
	"log"

	"github.com/williepotgieter/nft-maker/pkg/domain/models"
)

func (a *dbadapter) GetAllUsers() ([]models.User, error) {
	var (
		user  models.User
		users = []models.User{}
		rows  *sql.Rows
		err   error
	)

	rows, err = a.statements[GET_ALL_USERS].Query()
	if err != nil {
		log.Println(err)
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	for rows.Next() {
		user = models.User{}
		err = rows.Scan(&user.UUID, &user.Name, &user.Surname, &user.Email, &user.CreatedAt, &user.ModifiedAt)
		if err != nil {
			return []models.User{}, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (a *dbadapter) GetUser(uuid string) (models.User, error) {
	var (
		user models.User
		row  *sql.Row
		err  error
	)

	row = a.statements[GET_USER].QueryRow(uuid)

	err = row.Scan(&user.UUID, &user.Name, &user.Surname, &user.Email, &user.CreatedAt, &user.ModifiedAt)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
