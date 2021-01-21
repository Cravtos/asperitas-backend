package user

import (
	"github.com/cravtos/asperitas-backend/business/data/db"
)

func convertUserDBToInfo(usr db.FullUserDB) Info {
	return Info{
		ID:           usr.ID,
		Name:         usr.Name,
		PasswordHash: usr.PasswordHash,
		DateCreated:  usr.DateCreated,
	}
}
