package updating

func (s *dbservice) CloseConn() error {
	return s.db.CloseDBConn()
}

func (s *dbservice) Email(uuid string, email string) error {
	return s.db.UpdateEmail(uuid, email)
}
