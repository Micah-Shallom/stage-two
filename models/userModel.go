package models

import (
	"log"

	"gorm.io/gorm"
)

type User struct {
	UserID        string         `json:"userId" gorm:"primaryKey;unique;not null"`
	FirstName     string         `json:"firstName" gorm:"not null"`
	LastName      string         `json:"lastName" gorm:"not null"`
	Email         string         `json:"email" gorm:"unique;not null"`
	Password      string         `json:"password" gorm:"not null"`
	Phone         string         `json:"phone"`
	Organisations []Organisation `json:"organisations" gorm:"many2many:user_organisations;"` // many to many relationship
}

type UserModel struct {
	DB *gorm.DB
}

// Create a new user
func (m *UserModel) Create(user *User) error {
	return m.DB.Create(user).Error
}

// Get a user by email
func (m *UserModel) GetByEmail(email string) (*User, error) {
	var user User
	err := m.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Get a user by ID
func (m *UserModel) GetByID(id string) (*User, error) {
	var user User
	log.Printf("Querying user with ID: %s", id)
	err := m.DB.Where("user_id = ?", id).First(&user).Error
	if err != nil {
		log.Printf("Error fetching user by ID: %v", err)
		return nil, err
	}
	log.Printf("User found: %+v", user)
	return &user, nil
}
