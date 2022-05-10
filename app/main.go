package main

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

const (
	APP_NAME    string = "NFT Maker v0.0.1"
	DB_FILENAME string = "nftmaker.db"
	PORT        uint16 = 3000
)

var (
	ALGOD_ADDRESS, PS_TOKEN_KEY, PS_TOKEN, TESTNET_SOURCE_ACCOUNT string
	//go:embed secrets/api.yaml
	API_SECRETS_FILE []byte
)

// Read credentials from embedded api.yaml file
func init() {
	var (
		data = make(map[string]string)
		err  error
	)

	err = yaml.Unmarshal(API_SECRETS_FILE, &data)
	if err != nil {
		panic(err)
	}

	ALGOD_ADDRESS = data["algodAddress"]
	PS_TOKEN_KEY = data["psTokenKey"]
	PS_TOKEN = data["psToken"]
	TESTNET_SOURCE_ACCOUNT = data["testnetSourceAcc"]
}

func main() {
	db := NewDBConn(DB_FILENAME)
	bc := NewBlockchainClient(ALGOD_ADDRESS, PS_TOKEN_KEY, PS_TOKEN, TESTNET_SOURCE_ACCOUNT)
	api := NewRestAPI(APP_NAME, db, bc)

	api.Run(PORT)
}
