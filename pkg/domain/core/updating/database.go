package updating

func (s *dbservice) CloseConn() error {
	return s.db.CloseDBConn()
}
