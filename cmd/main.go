package main

import (
	"context"
	json "encoding/json"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
)

// const algodAddress = "Your Address"
// const algodToken = "Your Token"

const algodAddress = "https://testnet-algorand.api.purestake.io/idx2"
const algodToken = "Z2Q5c1rk3X2abzntNX53F6Gm9xJ9SPiG3MJjKgLu"

//var txHeaders = append([]*algod.Header{}, &algod.Header{"Content-Type", "application/json"})

// PrettyPrint prints Go structs
func PrettyPrint(data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}

func main() {
	account1 := crypto.GenerateAccount()
	account2 := crypto.GenerateAccount()
	account3 := crypto.GenerateAccount()
	address1 := account1.Address.String()
	address2 := account2.Address.String()
	address3 := account3.Address.String()

	mnemonic1, err := mnemonic.FromPrivateKey(account1.PrivateKey)
	if err != nil {
		return
	}
	mnemonic2, err := mnemonic.FromPrivateKey(account2.PrivateKey)
	if err != nil {
		return
	}
	mnemonic3, err := mnemonic.FromPrivateKey(account3.PrivateKey)
	if err != nil {
		return
	}
	fmt.Printf("1 : \"%s\"\n", address1)
	fmt.Printf("2 : \"%s\"\n", address2)
	fmt.Printf("3 : \"%s\"\n", address3)
	fmt.Printf("")
	fmt.Printf("Copy off accounts above and add TestNet Algo funds using the TestNet Dispenser at https://bank.testnet.algorand.network/\n")
	fmt.Printf("Copy off the following mnemonic code for future tutorial use\n")
	fmt.Printf("\n")
	fmt.Printf("mnemonic1 := \"%s\"\n", mnemonic1)
	fmt.Printf("mnemonic2 := \"%s\"\n", mnemonic2)
	fmt.Printf("mnemonic3 := \"%s\"\n", mnemonic3)

	// Initialize an algodClient
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		return
	}
	act, err := algodClient.AccountInformation(account1.Address.String()).Do(context.Background())
	if err != nil {
		fmt.Printf("failed to get account information: %s\n", err)
		return
	}
	fmt.Print("Account 1: ")
	PrettyPrint(act)
	act, err = algodClient.AccountInformation(account2.Address.String()).Do(context.Background())
	if err != nil {
		fmt.Printf("failed to get account information: %s\n", err)
		return
	}
	fmt.Print("Account 2: ")
	PrettyPrint(act)
	act, err = algodClient.AccountInformation(account3.Address.String()).Do(context.Background())
	if err != nil {
		fmt.Printf("failed to get account information: %s\n", err)
		return
	}
	fmt.Print("Account 3: ")
	PrettyPrint(act)
}
