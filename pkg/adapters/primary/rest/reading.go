package rest

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/williepotgieter/nft-maker/pkg/domain/models"
)

func (a *restadapter) handleGetAllUsers(c *fiber.Ctx) error {
	var (
		users []models.User
		err   error
	)

	users, err = a.db.Read.AllUsers()
	if err != nil {
		return generateResponse(c, fiber.StatusInternalServerError, "unable to get users")
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func (a *restadapter) handleGetUser(c *fiber.Ctx) error {
	var (
		user models.User
		uuid = c.Params("uuid")
		err  error
	)

	user, err = a.db.Read.User(uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return generateResponse(c, fiber.StatusNotFound, "user not found")
		}
		return generateResponse(c, fiber.StatusInternalServerError, "unable to get user")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}
