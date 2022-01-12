package database

type dbadapter struct{}

func NewDatabaseAdapter() *dbadapter {

	return &dbadapter{}
}
