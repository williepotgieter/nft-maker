package main

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"log"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

type Blockchain struct {
	client    *algod.Client
	sourceAcc string
}

type AlgoTransaction struct {
	FromAccount string `json:"from_account"`
	ToAccount   string `json:"to_account"`
	DecodedNote string `json:"decoded_note"`
	AmountSent  uint64 `json:"amount_sent"`
	Fee         uint64 `json:"fee"`
}

func NewBlockchainClient(algodAddress, psTokenKey, psToken, testnetSourceAcc string) *Blockchain {
	var (
		commonClient *common.Client
		algodClient  *algod.Client
		nodeStatus   models.NodeStatus
		err          error
	)

	commonClient, err = common.MakeClient(algodAddress, psTokenKey, psToken)
	if err != nil {
		panic(fmt.Errorf("failed to make common Algod client: %s\n", err))
	}

	algodClient = (*algod.Client)(commonClient)

	log.Println("Testing connection to the Algorand blockchain...")
	nodeStatus, err = algodClient.Status().Do(context.Background())
	if err != nil {
		panic(fmt.Errorf("error getting algod status: %s\n", err))
	}
	log.Printf("algod address: %s\n", algodAddress)
	log.Printf("algod last round: %d\n", nodeStatus.LastRound)
	log.Printf("algod time since last round: %d\n", nodeStatus.TimeSinceLastRound)
	log.Printf("algod catchup: %d\n", nodeStatus.CatchupTime)

	return &Blockchain{algodClient, testnetSourceAcc}
}

func CreateAlgorandAccount() (address, passphrase string, privatekey ed25519.PrivateKey, err error) {
	account := crypto.GenerateAccount()

	privatekey = account.PrivateKey

	passphrase, err = mnemonic.FromPrivateKey(privatekey)

	address = account.Address.String()

	return
}

func (bc *Blockchain) CheckAccountBalance(address string) (balance uint64, err error) {
	var accountInfo models.Account

	accountInfo, err = bc.client.AccountInformation(address).Do(context.Background())
	if err != nil {
		log.Printf("error getting Algorand account info: %s\n", err)
		return
	}

	balance = accountInfo.Amount

	return
}

func (bc *Blockchain) BuildAlgoTransferTxn(amount uint64, fromAddr, toAddr, txNote string) (txn types.Transaction, err error) {
	var (
		txParams                                types.SuggestedParams
		genID                                   string
		note, genHash                           []byte
		minFee, firstValidRound, lastValidRound uint64
	)

	txParams, err = bc.client.SuggestedParams().Do(context.Background())
	if err != nil {
		return
	}

	minFee = transaction.MinTxnFee
	note = []byte(txNote)
	genID = txParams.GenesisID
	genHash = txParams.GenesisHash
	firstValidRound = uint64(txParams.FirstRoundValid)
	lastValidRound = uint64(txParams.LastRoundValid)

	txn, err = transaction.MakePaymentTxnWithFlatFee(
		fromAddr,
		toAddr,
		minFee,
		amount,
		firstValidRound,
		lastValidRound,
		note,
		"",
		genID,
		genHash)
	if err != nil {
		return
	}

	return
}

func (bc *Blockchain) SignTxn(privateKey ed25519.PrivateKey, txn types.Transaction) (txID string, signedTxn []byte, err error) {
	txID, signedTxn, err = crypto.SignTransaction(privateKey, txn)
	if err != nil {
		return
	}

	return
}

func (bc *Blockchain) SubmitTxn(txID string, signedTxn []byte) (txInfo AlgoTransaction, err error) {
	var (
		sendResponse string
		confirmedTxn models.PendingTransactionInfoResponse
	)

	sendResponse, err = bc.client.SendRawTransaction(signedTxn).Do(context.Background())
	if err != nil {
		return
	}

	log.Printf("Submitted transaction %s\n", sendResponse)

	confirmedTxn, err = future.WaitForConfirmation(bc.client, txID, 4, context.TODO())
	if err != nil {
		return
	}

	txInfo = AlgoTransaction{
		FromAccount: confirmedTxn.Transaction.Txn.Sender.String(),
		ToAccount:   confirmedTxn.Transaction.Txn.Receiver.String(),
		DecodedNote: string(confirmedTxn.Transaction.Txn.Note),
		AmountSent:  uint64(confirmedTxn.Transaction.Txn.Amount),
		Fee:         uint64(confirmedTxn.Transaction.Txn.Fee),
	}

	return
}
