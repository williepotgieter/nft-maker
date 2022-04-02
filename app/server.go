package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type RestApi struct {
	*DBConn
	app      *fiber.App
	validate *validator.Validate
}

func NewRestAPI(name string, c *DBConn) *RestApi {
	app := fiber.New(fiber.Config{
		AppName: name,
	})

	// Add middleware
	app.Use(cors.New())
	app.Use(logger.New())

	v := validator.New()

	return &RestApi{c, app, v}
}

func (s *RestApi) Run(port uint16) {
	// Add routes
	s.setupV1Routes()

	s.app.Listen(fmt.Sprintf("localhost:%d", port))
}

func (s *RestApi) setupV1Routes() {
	// Group all V1 endpoints
	v1 := s.app.Group("/api/v1")

	// Users endpoints
	users := v1.Group("/users")
	users.Post("/register", s.HandleRegisterNewUser)
	users.Get("/", s.HandleGetAllUsers)
	users.Get("/:id", s.HandleGetUser)
	users.Patch("/:id/password", s.HandleUpdateUserPassword)
	users.Delete("/:id", s.HandleDeleteUser)
}
