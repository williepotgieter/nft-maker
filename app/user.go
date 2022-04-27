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
	Accounts []Account `gorm:"constraint:OnDelete:CASCADE;"`
}
type UserResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
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

// HandleRegisterNewUser godoc
// @Summary      Create a new user
// @Description  Register a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        Body    body     UserRegistration  true  "New user details"
// @Success      201  {object} HTTPResp
// @Failure      400  {object}  HTTPResp
// @Failure      409  {object}  HTTPResp
// @Failure      500  {object}  HTTPResp
// @Router       /users/register [post]
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

// HandleGetAllUsers godoc
// @Summary      List all users
// @Description  Lists all users stored in the database
// @Tags         users
// @Produce      json
// @Success      200  {array} UserResponse
// @Failure      500  {object}  HTTPResp
// @Router       /users [get]
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
			CreatedAt: user.CreatedAt.UTC(),
			UpdatedAt: user.UpdatedAt.UTC(),
			Name:      user.Name,
			Email:     user.Email,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// HandleGetUser godoc
// @Summary      Specific user details
// @Description  Get details of a specific user
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {array} UserResponse
// @Failure      400  {object}  HTTPResp
// @Failure      404  {object}  HTTPResp
// @Failure      500  {object}  HTTPResp
// @Router       /users/{id} [get]
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
		CreatedAt: user.CreatedAt.UTC(),
		UpdatedAt: user.UpdatedAt.UTC(),
		Name:      user.Name,
		Email:     user.Email,
	})
}

// HandleUpdateUserPassword godoc
// @Summary      Update password
// @Description  Updates an existing user's password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param        Body    body     UserPassword  true  "New password"
// @Success      200  {array} UserResponse
// @Failure      400  {object}  HTTPResp
// @Failure      500  {object}  HTTPResp
// @Router       /users/{id}/password [patch]
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

// HandleDeleteUser godoc
// @Summary      Delete user
// @Description  Delete a specific user
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {array} UserResponse
// @Failure      400  {object}  HTTPResp
// @Failure      500  {object}  HTTPResp
// @Router       /users/{id} [delete]
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
