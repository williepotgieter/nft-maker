package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AccountInfo struct {
	UserID  uint   `json:"user_id"`
	Address string `json:"address"`
}

type Account struct {
	gorm.Model
	AccountInfo
	Passphrase string
	PrivateKey string
	UserID     uint
}

type AccountResponse struct {
	ID        uint      `json:"id"`
	Address   string    `json:"address"`
	Balance   uint64    `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateNewAccount CREATES a new Algorand account for a specific
// user and stores the account details in the database
func (c *DBConn) CreateNewAlgorandAccount(userID uint) (err error) {
	var (
		address, passphrase string
		privatekey          ed25519.PrivateKey
	)

	address, passphrase, privatekey, err = CreateAlgorandAccount()
	if err != nil {
		return err
	}

	err = c.db.Create(&Account{
		AccountInfo: AccountInfo{userID, address},
		Passphrase:  passphrase,
		PrivateKey:  base64.StdEncoding.EncodeToString(privatekey),
		UserID:      userID,
	}).Error

	return
}

// GetUserAccounts returns all accounts for a specific user
func (c *DBConn) GetUserAccounts(userID uint) (accounts []Account, err error) {
	err = c.db.Find(&accounts, "user_id = ?", userID).Error

	return
}

// HandleCreateNewAlgorandAccount godoc
// @Summary      Create a new account
// @Description  Create new Algorand account for a specific user
// @Tags         accounts
// @Produce      json
// @Param        userId   path      int  true  "User ID"
// @Success      201  {object} HTTPResp
// @Failure      400  {object}  HTTPResp
// @Failure      500  {object}  HTTPResp
// @Router       /accounts/{userId}/new [post]
func (s *RestApi) HandleCreateNewAlgorandAccount(c *fiber.Ctx) error {
	var (
		userId int
		err    error
	)

	userId, err = c.ParamsInt("userId")
	if err != nil {
		return httpResponse(c, fiber.StatusBadRequest, "invalid request params")
	}

	err = s.DBConn.CreateNewAlgorandAccount(uint(userId))
	if err != nil {
		log.Println(err.Error())
		return httpResponse(c, fiber.StatusInternalServerError, "unable to create new Algorand account")
	}

	return httpResponse(c, fiber.StatusCreated, "successfully created new Algorand account")
}

// HandleGetUserAccounts godoc
// @Summary      Get user's accounts
// @Description  List all accounts for a specific user
// @Tags         accounts
// @Produce      json
// @Param        userId   path      int  true  "User ID"
// @Success      200  {array} AccountResponse
// @Failure      400  {object}  HTTPResp
// @Failure      500  {object}  HTTPResp
// @Router       /accounts/{userId}/all [get]
func (s *RestApi) HandleGetUserAccounts(c *fiber.Ctx) error {
	var (
		userId   int
		accounts []Account
		response []AccountResponse = []AccountResponse{}
		balance  uint64
		err      error
	)

	userId, err = c.ParamsInt("userId")
	if err != nil {
		return httpResponse(c, fiber.StatusBadRequest, "invalid request params")
	}

	accounts, err = s.DBConn.GetUserAccounts(uint(userId))
	if err != nil {
		log.Println(err.Error())
		return httpResponse(c, fiber.StatusInternalServerError, "unable to query user's accounts")
	}

	for _, account := range accounts {
		balance, err = s.Blockchain.CheckAccountBalance(account.Address)
		if err != nil {
			return httpResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("error while trying to retrieve account balance for %s\n", account.Address))
		}
		response = append(response, AccountResponse{
			ID:        account.ID,
			Address:   account.Address,
			Balance:   balance,
			CreatedAt: account.CreatedAt.UTC(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
