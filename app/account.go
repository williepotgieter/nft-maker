package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	_ "github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TransferQuery struct {
	Amount uint64 `query:"amount"`
	ToAddr string `query:"to_addr"`
}

type TransferBody struct {
	TxNote string `json:"tx_note"`
}

type TxHistoryQuery struct {
	Address   string `query:"address"`
	StartTime string `query:"start_time"`
}

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

func (c *DBConn) GetUserAccountPrivateKey(address string) (privateKey ed25519.PrivateKey, err error) {
	var account Account
	err = c.db.Where("address = ?", address).First(&account).Error
	if err != nil {
		return
	}

	privateKey, err = base64.StdEncoding.DecodeString(account.PrivateKey)
	if err != nil {
		return
	}

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
		respCh   chan AccountResponse
		errCh    = make(chan error)
		accounts []Account
		response []AccountResponse = []AccountResponse{}
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

	// Make concurrent blockchain requests to get balances for all of the user's Algorand accounts
	respCh = make(chan AccountResponse, len(accounts))
	for _, account := range accounts {
		go func(accID uint, accAddress string, accCreatedAt time.Time) {
			var (
				accBal uint64
				reqErr error
			)

			accBal, reqErr = s.Blockchain.CheckAccountBalance(accAddress)
			if reqErr != nil {
				errCh <- fmt.Errorf("error while trying to retrieve account balance for %s\n", accAddress)
			}

			respCh <- AccountResponse{
				ID:        accID,
				Address:   accAddress,
				Balance:   accBal,
				CreatedAt: accCreatedAt,
			}
		}(account.ID, account.Address, account.CreatedAt.UTC())
	}

	for i := 0; i < len(accounts); i++ {
		select {
		case acc := <-respCh:
			response = append(response, acc)
		case errResp := <-errCh:
			close(respCh)
			close(errCh)
			return httpResponse(c, fiber.StatusInternalServerError, errResp.Error())
		case <-time.After(time.Second * 15):
			return httpResponse(c, fiber.StatusRequestTimeout, "request took too long")
		}
	}

	close(errCh)

	// Sort accounts by ID. This works because account ID's are created in sequence in the database
	sort.Slice(response, func(i, j int) bool { return response[i].ID < response[j].ID })

	return c.Status(fiber.StatusOK).JSON(response)
}

// HandleTransferAlgo godoc
// @Summary      Make an Algorand transfer
// @Description  Transfer Algorand from one account to another
// @Tags         accounts
// @Produce      json
// @Param        userId   path      int  true  "User ID"
// @Param        sourceAcc   path      string  true  "Source account"
// @Param        amount   query      uint64  true  "Transfer amount (in micro Algo)"
// @Param        to_addr   query      string  true  "Destination address"
// @Param        tx_note   body      TransferBody  true  "Transaction note"
// @Success      200  {object} AlgoTransaction
// @Failure      400  {object}  HTTPResp
// @Failure      500  {object}  HTTPResp
// @Router       /accounts/{userId}/{sourceAcc}/transfer [post]
func (s *RestApi) HandleTransferAlgo(c *fiber.Ctx) error {
	var (
		sourceAcc     = c.Params("sourceAcc")
		transferQuery = new(TransferQuery)
		transferBody  = new(TransferBody)
		privateKey    ed25519.PrivateKey
		err           error
	)

	err = c.QueryParser(transferQuery)
	if err != nil {
		return httpResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err = c.BodyParser(transferBody)
	if err != nil {
		return httpResponse(c, fiber.StatusBadRequest, err.Error())
	}

	privateKey, err = s.GetUserAccountPrivateKey(sourceAcc)
	if err != nil {
		return httpResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	txn, err := s.Blockchain.BuildAlgoTransferTxn(
		transferQuery.Amount,
		sourceAcc,
		transferQuery.ToAddr,
		transferBody.TxNote)
	if err != nil {
		return httpResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	txID, signedTxn, err := s.Blockchain.SignTxn(privateKey, txn)
	if err != nil {
		return httpResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	txInfo, err := s.Blockchain.SubmitTxn(txID, signedTxn)
	if err != nil {
		return httpResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(txInfo)
}

func (s *RestApi) HandleGetAccountHistory(c *fiber.Ctx) error {
	var (
		q       = new(TxHistoryQuery)
		history models.TransactionsResponse
		err     error
	)

	err = c.QueryParser(q)
	if err != nil {
		return httpResponse(c, fiber.StatusBadRequest, err.Error())
	}

	history, err = s.Blockchain.GetAccountHistory(q.Address, "2020-06-03T10:00:00-05:00")
	if err != nil {
		return httpResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(history)
}
