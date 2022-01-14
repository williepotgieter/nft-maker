package rest

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/williepotgieter/nft-maker/pkg/domain/models"
)

func (a *restadapter) handleCreateNewUser(c *fiber.Ctx) error {
	var (
		newUser = models.User{}
		err     error
	)

	//newUser := models.User{}

	err = json.Unmarshal(c.Body(), &newUser)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid request body")
	}

	err = a.db.Add.User(newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("unable to save user to database")
	}

	return c.Status(fiber.StatusCreated).SendString("New user added")
}
