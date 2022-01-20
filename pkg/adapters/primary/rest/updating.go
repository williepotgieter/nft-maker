package rest

import (
	"github.com/gofiber/fiber/v2"
)

type UpdateEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (a *restadapter) handleUpdateEmail(c *fiber.Ctx) error {
	var (
		uuid = c.Params("uuid")
		body = UpdateEmailRequest{}
		err  error
	)

	err = a.validate.Var(uuid, "required,uuid")
	if err != nil {
		return generateResponse(c, fiber.StatusBadRequest, "invalid request params")
	}

	err = c.BodyParser(&body)
	if err != nil {
		return generateResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	err = a.validate.Struct(body)
	if err != nil {
		return generateResponse(c, fiber.StatusBadRequest, "invalid email")
	}

	err = a.db.Update.Email(uuid, body.Email)
	if err != nil {
		if err.Error() == "notfound" {
			return generateResponse(c, fiber.StatusNotFound, "unable to update email")
		}
		return generateResponse(c, fiber.StatusInternalServerError, "unable to update email")
	}

	return generateResponse(c, fiber.StatusOK, "user email updated")
}
