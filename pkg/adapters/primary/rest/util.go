package rest

import "github.com/gofiber/fiber/v2"

type stdResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// return c.Status(fiber.StatusInternalServerError).SendString("unable to save user to database")

func generateResponse(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(stdResponse{code, message})
}
