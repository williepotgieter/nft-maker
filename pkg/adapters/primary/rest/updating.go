package rest

import (
	"github.com/gofiber/fiber/v2"
)

type UpdateEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password" validate:"required,min=8,max=32"`
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

func (a *restadapter) handleUpdatePassword(c *fiber.Ctx) error {
	var (
		uuid = c.Params("uuid")
		body = UpdatePasswordRequest{}
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
		return generateResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	err = a.db.Update.Password(uuid, body.Password)
	if err != nil {
		if err.Error() == "notfound" {
			return generateResponse(c, fiber.StatusNotFound, "unable to update password")
		}
		return generateResponse(c, fiber.StatusInternalServerError, "unable to update password")
	}

	return generateResponse(c, fiber.StatusOK, "user password updated")
}
