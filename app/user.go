package main

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SaveNewUser CREATES a new user in the database
func (c *DBConn) SaveNewUser(name, email, password string) (err error) {

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return
	}

	c.db.Create(&User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	})
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

// UpdateUserPassword UPDATES a specific user's password
func (c *DBConn) UpdateUserPassword(id uint, password string) (err error) {
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
func (c *DBConn) DeleteUser(id string) (err error) {
	err = c.db.Delete(&User{}, id).Error

	return
}
