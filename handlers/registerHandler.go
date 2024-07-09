package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Micah-Shallom/stage-two/models"
	"github.com/Micah-Shallom/stage-two/utils"
	"github.com/Micah-Shallom/stage-two/validator"
	"github.com/gin-gonic/gin"
)

// registerUserHandler is a handler function to register a new user
func (h *Handlers) RegisterUserHandler(c *gin.Context) {

	// Create an anonymous struct to hold the data that we expect to be in the request body
	var req validator.RegisterReq

	// Parse the req body into the anonymous struct
	err := utils.ReadRequest(c, &req)
	if err != nil {
		utils.BadRequestResponse(c, "Registration unsuccessful", 400, err)
		return
	}

	//perform validation for the request struct
	errorResponse, err := req.Validate()

	if errorResponse != nil || err != nil {
		utils.ValidationErrorResponse(c, json.RawMessage(errorResponse))
		return
	}

	//get hashed password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.LogError(c, err)
		// app.errorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create a new User object containing the data that we read from the request body
	user := &models.User{
		UserID:    utils.GenerateUUID(),
		FirstName: req.Firstname,
		LastName:  req.Lastname,
		Email:     req.Email,
		Password:  hashedPassword,
		Phone:     req.Phone,
	}

	//save user model to the database
	err = h.App.Models.Users.Create(user)
	if err != nil {
		utils.BadRequestResponse(c, "Registration unsuccessful", 422, err)
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
		utils.BadRequestResponse(c, "Registration unsuccessful", 400, err)
		return
	}

	//generate jwt token
	accessToken, err := utils.GenerateJWT(user)
	if err != nil {
		utils.LogError(c, err)
		utils.BadRequestResponse(c, "Registration unsuccessful", 400, err)
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
