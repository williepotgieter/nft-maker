package main

import "github.com/gofiber/fiber/v2"

type HTTPResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func httpResponse(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(HTTPResp{code, message})
}
