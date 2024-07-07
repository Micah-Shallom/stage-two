package handlers

import (
	"net/http"

	"github.com/Micah-Shallom/stage-two/models"
	"github.com/Micah-Shallom/stage-two/utils"
	"github.com/gin-gonic/gin"
)

// registerUserHandler is a handler function to register a new user
func (h *Handlers) RegisterUserHandler(c *gin.Context) {

	// Create an anonymous struct to hold the data that we expect to be in the request body
	var input struct {
		FirstName string `json:"firstName" validate:"required,alpha"`
		LastName  string `json:"lastName" validate:"required,alpha"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required,min=8"`
		Phone     string `json:"phone" validate:"required"`
	}

	// Parse the request body into the anonymous struct
	valid, errors := utils.ReadRequest(c, &input)

	if !valid {
		utils.ValidationErrorResponse(c, errors)
		return
	}

	//get hashed password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.LogError(c, err)
		// app.errorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create a new User object containing the data that we read from the request body
	user := &models.User{
		UserID:    utils.GenerateUUID(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  hashedPassword,
		Phone:     input.Phone,
	}

	//save user model to the database
	err = h.App.Models.Users.Create(user)
	if err != nil {
		utils.BadRequestResponse(c, "Registration unsuccessful", http.StatusBadRequest, err)
		return
	}

	//create a default organization for the user
	org := &models.Organisation{
		OrgID:       utils.GenerateUUID(),
		Name:        user.FirstName + "'s Organisation",
		Description: "Default organisation for " + user.FirstName,
		Users:       []models.User{*user},
	}

	//save organisation model to the database
	err = h.App.Models.Organisations.Create(org)
	if err != nil {
		utils.LogError(c, err)
		utils.BadRequestResponse(c, "Registration unsuccessful", http.StatusBadRequest, err)
		return
	}

	//generate jwt token
	accessToken, err := utils.GenerateJWT(user)
	if err != nil {
		utils.LogError(c, err)
		utils.BadRequestResponse(c, "Registration unsuccessful", http.StatusBadRequest, err)
		return
	}

	//Respond with a 201 Created status code
	response := utils.Envelope{
		"status":  "success",
		"message": "Registration successful",
		"data": map[string]any{
			"accessToken": accessToken,
			"user": map[string]any{
				"userId":    user.UserID,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"email":     user.Email,
				"phone":     user.Phone,
			},
		},
	}
	utils.WriteResponse(c, http.StatusCreated, response, nil)
}
