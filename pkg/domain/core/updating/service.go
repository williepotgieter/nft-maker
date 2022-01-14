package updating

// Database service
type Database interface {
	CloseDBConn() error
}

type DatabaseService interface {
	CloseConn() error
}

type dbservice struct {
	db Database
}

func NewDBService(db Database) DatabaseService {
	return &dbservice{db}
}
