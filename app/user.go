package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserPassword struct {
	Password string `json:"password" validate:"required"`
}
type UserRegistration struct {
	Name  string `json:"name" validate:"required,min=1"`
	Email string `json:"email" gorm:"unique" validate:"required,email"`
	UserPassword
}
type User struct {
	gorm.Model
	UserRegistration
}

type UserResponse struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

// SaveNewUser CREATES a new user in the database
func (c *DBConn) SaveNewUser(newUser UserRegistration) (err error) {

	hashedPassword, err := HashPassword(newUser.Password)
	if err != nil {
		return
	}

	err = c.db.Create(&User{
		UserRegistration: UserRegistration{
			Name:  newUser.Name,
			Email: newUser.Email,
			UserPassword: UserPassword{
				Password: hashedPassword,
			},
		},
	}).Error
	return
}

// GetUser READS a specific user from the database
func (c *DBConn) GetUser(id uint) (user User, err error) {
	err = c.db.First(&user, id).Error

	return
}

// GetAllUsers READS all users from the database
func (c *DBConn) GetAllUsers() (users []User, err error) {
	err = c.db.Find(&users).Error

	return
}

// UserPassword UPDATES a specific user's password
func (c *DBConn) UserPassword(id uint, password string) (err error) {
	var (
		user           User
		hashedPassword string
	)

	hashedPassword, err = HashPassword(password)

	err = c.db.First(&user, id).Error
	if err != nil {
		return
	}

	err = c.db.Model(&user).Update("Password", hashedPassword).Error

	return
}

// DeleteUser DELETES a specific user from the database
func (c *DBConn) DeleteUser(id uint) (err error) {
	err = c.db.Unscoped().Delete(&User{}, id).Error

	return
}

// HandleRegisterNewUser is a http handler function
// used to save a new user to the database
func (s *RestApi) HandleRegisterNewUser(c *fiber.Ctx) error {
	var (
		user = UserRegistration{}
		err  error
	)

	err = c.BodyParser(&user)
	if err != nil {
		return httpResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	err = s.validate.Struct(user)
	if err != nil {
		return httpResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	err = s.DBConn.SaveNewUser(user)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return httpResponse(c, fiber.StatusConflict, "unable to save user")
		}
		return httpResponse(c, fiber.StatusInternalServerError, "unable to save user")
	}

	return httpResponse(c, fiber.StatusCreated, "successfully registered new user")
}

// HandleGetAllUsers is a http handler function that
// returns a list of all users stored in the database
func (s *RestApi) HandleGetAllUsers(c *fiber.Ctx) error {
	var (
		users    []User
		response = []UserResponse{}
		err      error
	)

	users, err = s.DBConn.GetAllUsers()
	if err != nil {
		return httpResponse(c, fiber.StatusInternalServerError, "unable to read users from database")
	}

	for _, user := range users {
		response = append(response, UserResponse{
			ID:        user.ID,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
			Name:      user.Name,
			Email:     user.Email,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// HandleGetUser is a http handler function that
// returns the details of a specific user
func (s *RestApi) HandleGetUser(c *fiber.Ctx) error {
	var (
		id   int
		user User
		err  error
	)

	id, err = c.ParamsInt("id")
	if err != nil {
		return httpResponse(c, fiber.StatusBadRequest, "invalid request params")
	}

	user, err = s.DBConn.GetUser(uint(id))
	if err != nil {
		return httpResponse(c, fiber.StatusNotFound, "unknown user")
	}

	return c.Status(fiber.StatusOK).JSON(UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Name:      user.Name,
		Email:     user.Email,
	})
}

// HandleUpdateUserPassword is a http handler function that
// updates an existing user's password
func (s *RestApi) HandleUpdateUserPassword(c *fiber.Ctx) error {
	var (
		id             int
		hashedPassword string
		user           UserPassword
		err            error
	)

	id, err = c.ParamsInt("id")
	if err != nil {
		return httpResponse(c, fiber.StatusBadRequest, "invalid request params")
	}

	err = c.BodyParser(&user)
	if err != nil {
		return httpResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	err = s.validate.Struct(user)
	if err != nil {
		return httpResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	hashedPassword, err = HashPassword(user.Password)
	if err != nil {
		return httpResponse(c, fiber.StatusInternalServerError, "unable to update user password")
	}

	err = s.DBConn.UserPassword(uint(id), hashedPassword)
	if err != nil {
		return httpResponse(c, fiber.StatusInternalServerError, "unable to update user password")
	}

	return httpResponse(c, fiber.StatusOK, "successfully updated user password")
}

// HandleDeleteUser is a http handler function that
// deletes a specific user
func (s *RestApi) HandleDeleteUser(c *fiber.Ctx) error {
	var (
		id  int
		err error
	)

	id, err = c.ParamsInt("id")
	if err != nil {
		return httpResponse(c, fiber.StatusBadRequest, "invalid request params")
	}

	err = s.DBConn.DeleteUser(uint(id))
	if err != nil {
		return httpResponse(c, fiber.StatusInternalServerError, "unable to delete user")
	}

	return httpResponse(c, fiber.StatusOK, "user deleted")
}
