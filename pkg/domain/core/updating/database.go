package updating

import (
	"errors"

	"github.com/williepotgieter/nft-maker/pkg/domain/utils"
)

func (s *dbservice) CloseConn() error {
	return s.db.CloseDBConn()
}

func (s *dbservice) Email(uuid string, email string) error {
	return s.db.UpdateEmail(uuid, email)
}

func (s *dbservice) Password(uuid string, password string) error {
	var (
		hashedPassword string
		err            error
	)

	hashedPassword, err = utils.HashPassword(password)
	if err != nil {
		return errors.New("unable to hash password")
	}

	return s.db.UpdatePassword(uuid, hashedPassword)
}
