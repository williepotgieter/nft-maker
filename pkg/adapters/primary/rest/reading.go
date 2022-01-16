package rest

import (
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
