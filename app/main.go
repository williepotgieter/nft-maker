package main

const (
	APP_NAME    string = "NFT Maker v0.0.1"
	DB_FILENAME string = "nftmaker.db"
	PORT        uint16 = 3000
)

func main() {
	db := NewDBConn(DB_FILENAME)
	api := NewRestAPI(APP_NAME, db)

	api.Run(PORT)
}
