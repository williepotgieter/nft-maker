package updating

// Database service
type Database interface {
	CloseDBConn() error
	UpdateEmail(uuid string, email string) error
}

type DatabaseService interface {
	CloseConn() error
	Email(uuid string, email string) error
}

type dbservice struct {
	db Database
}

func NewDBService(db Database) DatabaseService {
	return &dbservice{db}
}
