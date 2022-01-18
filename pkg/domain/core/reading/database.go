package reading

import "github.com/williepotgieter/nft-maker/pkg/domain/models"

func (s *dbservice) AllUsers() ([]models.User, error) {
	return s.db.GetAllUsers()
}

func (s *dbservice) User(uuid string) (models.User, error) {
	return s.db.GetUser(uuid)
}
