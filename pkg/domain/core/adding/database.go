package adding

import (
	"time"

	"github.com/google/uuid"
	"github.com/williepotgieter/nft-maker/pkg/domain/models"
	"github.com/williepotgieter/nft-maker/pkg/domain/utils"
)

func (s *dbservice) User(user models.User) error {
	var (
		id        uuid.UUID
		timestamp int64
		err       error
	)

	id, err = uuid.NewUUID()
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	timestamp = time.Now().Unix()

	return s.db.CreateUser(models.User{
		CreatedAt:  timestamp,
		ModifiedAt: timestamp,
		UUID:       id.String(),
		Name:       user.Name,
		Surname:    user.Surname,
		Email:      user.Email,
		Password:   hashedPassword,
	})
}
