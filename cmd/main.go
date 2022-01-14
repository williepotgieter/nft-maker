package main

import (
	"github.com/williepotgieter/nft-maker/pkg/adapters/primary/rest"
	"github.com/williepotgieter/nft-maker/pkg/adapters/secondary/database"
	"github.com/williepotgieter/nft-maker/pkg/domain/ports/repository"
)

func main() {
	// Database adapter
	dbAdapter := database.NewDBAdapter()
	dbPort := repository.NewDatabasePort(dbAdapter)

	restService := rest.NewRESTAdapter(dbPort)

	restService.Run()
}
