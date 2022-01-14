package database

func (a *dbadapter) CloseDBConn() error {
	return a.db.Close()
}
