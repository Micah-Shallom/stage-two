package models

import "gorm.io/gorm"

type Organisation struct {
	OrgID       string `json:"orgId" gorm:"primaryKey;unique;not null"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	Users       []User `json:"users" gorm:"many2many:user_organisations;foreignKey:OrgID;joinForeignKey:org_id;References:UserID;joinReferences:user_id"`
}
type OrganisationModel struct {
	DB *gorm.DB
}

// Create a new organisation
func (m *OrganisationModel) Create(organisation *Organisation) error {
	return m.DB.Create(organisation).Error
}

// Get Organizations by user ID
func (m *OrganisationModel) GetByUserID(id string) ([]Organisation, error) {
	var orgs []Organisation

	// Adjust the query to use the correct column names from user_organisations table
	err := m.DB.
		Preload("Users").
		Joins("JOIN user_organisations ON user_organisations.org_id = organisations.org_id").
		Joins("JOIN users ON user_organisations.user_id = users.user_id").
		Where("users.user_id = ?", id).
		Find(&orgs).Error
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

//Get Organisation by org id
func (m *OrganisationModel) GetByOrgID(id string) (*Organisation, error) {
	var org Organisation
	// err := m.DB.Preload("Users").Where("org_id = ?", id).First(&org).Error
	err := m.DB.Where("org_id = ?", id).First(&org).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (m *OrganisationModel) IsUserInOrganisation(orgID, userID string) (bool, error) {
	var org Organisation
	err := m.DB.
		Where("org_id = ?", orgID).
		Preload("Users").
		First(&org).Error
	if err != nil {
		return false, err
	}

	for _, user := range org.Users {
		if user.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}

func (m *OrganisationModel) AddUserToOrganisation(orgID, userID string) error {
	// Add user to organisation
	err := m.DB.Exec("INSERT INTO user_organisations (org_id, user_id) VALUES (?, ?)", orgID, userID).Error
	if err != nil {
		return err
	}
	return nil
}
