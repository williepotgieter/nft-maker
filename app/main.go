package main

import "fmt"

var nftmaker *DBConn

func init() {
	nftmaker := NewDBConn("nftmaker.db")
}

func main() {
	fmt.Println("Hello, World!")
}
