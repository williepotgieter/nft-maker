package rest

import (
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/williepotgieter/nft-maker/pkg/domain/ports/repository"
)

type restadapter struct {
	app      *fiber.App
	db       *repository.DatabasePort
	validate *validator.Validate
}

func NewRESTAdapter(db *repository.DatabasePort) *restadapter {
	app := fiber.New(fiber.Config{
		AppName: "NFT Maker API v0.0.1",
	})

	// Add middleware
	app.Use(cors.New(cors.ConfigDefault))
	app.Use(logger.New())

	return &restadapter{
		app:      app,
		db:       db,
		validate: validator.New(),
	}
}

func (a *restadapter) InitV1Routes() {
	// API V1
	v1 := a.app.Group("/v1")

	// Groups
	users := v1.Group("/users")

	// Users endpoints
	users.Post("/", a.handleCreateNewUser)
	users.Get("/", a.handleGetAllUsers)
	users.Get("/:uuid", a.handleGetUser)
	users.Patch("/:uuid/email", a.handleUpdateEmail)
	users.Patch("/:uuid/password", a.handleUpdatePassword)
}

func (a *restadapter) Run() {
	var (
		port string
		err  error
	)

	port = fmt.Sprintf(":%s", os.Getenv("PORT"))

	defer func() {
		err = a.db.Update.CloseConn()
		if err != nil {
			log.Println(err)
		}
	}()

	a.InitV1Routes()

	a.app.Listen(port)
}
