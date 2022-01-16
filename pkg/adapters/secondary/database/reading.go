package database

import (
	"log"

	"github.com/williepotgieter/nft-maker/pkg/domain/models"
)

func (a *dbadapter) GetAllUsers() ([]models.User, error) {
	var (
		user  models.User
		users = []models.User{}
		err   error
	)

	rows, err := a.statements[GET_ALL_USERS].Query()
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
