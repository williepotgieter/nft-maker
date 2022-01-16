package rest

import (
	"encoding/json"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/williepotgieter/nft-maker/pkg/domain/models"
)

const MYSQL_KEY_EXITS = 1062

func (a *restadapter) handleCreateNewUser(c *fiber.Ctx) error {
	var (
		newUser = models.User{}
		err     error
	)

	err = json.Unmarshal(c.Body(), &newUser)
	if err != nil {
		return generateResponse(c, fiber.StatusBadRequest, "invalid request")
	}

	err = a.db.Add.User(newUser)
	if err != nil {
		if err.(*mysql.MySQLError).Number == MYSQL_KEY_EXITS {
			return generateResponse(c, fiber.StatusConflict, "user already exists")
		}
		return generateResponse(c, fiber.StatusInternalServerError, "unable to save user to database")
	}

	return generateResponse(c, fiber.StatusCreated, "new user added")
}
