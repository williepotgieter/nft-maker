package main

import "github.com/gofiber/fiber/v2"

func httpResponse(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{
		"code":    code,
		"message": message,
	})
}
