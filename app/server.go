package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	docs "nft-maker/docs"

	swagger "github.com/arsmn/fiber-swagger/v2"
)

type RestApi struct {
	*DBConn
	*Blockchain
	app      *fiber.App
	validate *validator.Validate
}

func NewRestAPI(name string, c *DBConn, a *Blockchain) *RestApi {
	app := fiber.New(fiber.Config{
		AppName: name,
	})

	// Add middleware
	app.Use(cors.New())
	app.Use(logger.New())

	v := validator.New()

	return &RestApi{c, a, app, v}
}

// @title NFT Maker API
// @version 0.0.1
// @description Simple API to demonstrate the creation of NFT's on the Algorand blockchain, stored on ipfs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1

func (s *RestApi) Run(port uint16) {
	// Add routes
	s.setupV1Routes()

	// Add swagger
	s.setupSwagger(port)

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

	// Accounts endpoints
	accounts := v1.Group("/accounts")
	accounts.Post("/:userId/new", s.HandleCreateNewAlgorandAccount)
	accounts.Post("/:userId/:accAddress/transfer", s.HandleTransferAlgo)
	accounts.Get("/:userId/all", s.HandleGetUserAccounts)
}

func (s *RestApi) setupSwagger(port uint16) {
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", port)
	s.app.Get("/swagger/*", swagger.HandlerDefault)
}
