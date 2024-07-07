package models

import "gorm.io/gorm"

type Models struct {
	Users         UserModel
	Organisations OrganisationModel
}

func NewModels(db *gorm.DB) Models {
	return Models{
		Users: UserModel{
			DB: db,
		},
		Organisations: OrganisationModel{
			DB: db,
		},
	}
}
