package persistence

import (
	"database/sql"
	models "postapi/internal/domain"

	"github.com/jmoiron/sqlx"
)

type ProfileRepositoryImpl struct {
	db *sqlx.DB
}

func (pR *ProfileRepositoryImpl) FindByUsername(username string) (*models.Profile, error) {
	profile := &models.Profile{}
	err := pR.db.Get(profile, getProfileSchema, username)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (pR *ProfileRepositoryImpl) Create(p *models.Profile) error {
	var profilePicture string
	if p.ProfilePicture == "" {
		profilePicture = "https://i.redd.it/j6mkb6p73h791.jpg"
	} else {
		profilePicture = p.ProfilePicture
	}
	_, err := pR.db.Exec(insertProfileSchema, p.Username, p.Description, profilePicture)
	return err
}
func (pR *ProfileRepositoryImpl) Update(p *models.Profile) error {
	result, err := pR.db.Exec(updateProfileSchema, p.Username, p.Description, p.ProfilePicture)
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
