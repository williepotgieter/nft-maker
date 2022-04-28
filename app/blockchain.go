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
	"github.com/algorand/go-algorand-sdk/mnemonic"
)

type Blockchain struct {
	client *algod.Client
}

func NewBlockchainClient(algodAddress, psTokenKey, psToken string) *Blockchain {
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

	return &Blockchain{algodClient}
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
