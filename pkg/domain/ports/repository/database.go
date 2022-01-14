package repository

import (
	"github.com/williepotgieter/nft-maker/pkg/domain/core/adding"
	"github.com/williepotgieter/nft-maker/pkg/domain/core/updating"
)

type Database interface {
	adding.Database
	updating.Database
}

type DatabasePort struct {
	Add    adding.DatabaseService
	Update updating.DatabaseService
}

func NewDatabasePort(db Database) *DatabasePort {
	return &DatabasePort{
		Add:    adding.NewDBService(db),
		Update: updating.NewDBService(db),
	}
}
