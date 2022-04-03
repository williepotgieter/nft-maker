package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	_ "nft-maker/docs"

	swagger "github.com/arsmn/fiber-swagger/v2"
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

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /

func (s *RestApi) Run(port uint16) {
	// Add routes
	s.setupV1Routes()

	// Add swagger
	s.setupSwagger()

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

func (s *RestApi) setupSwagger() {
	s.app.Get("/swagger/*", swagger.HandlerDefault)
}
