package main

import (
	"crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
)

func CreateAlgorandAccount() (address, passphrase string, privatekey ed25519.PrivateKey, err error) {
	account := crypto.GenerateAccount()

	privatekey = account.PrivateKey

	passphrase, err = mnemonic.FromPrivateKey(privatekey)

	address = account.Address.String()

	return
}
