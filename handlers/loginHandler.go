package handlers

import (
	"encoding/json"

	"github.com/Micah-Shallom/stage-two/utils"
	"github.com/Micah-Shallom/stage-two/validator"
	"github.com/gin-gonic/gin"
)

// loginHandler is a handler function to login a user
func (h *Handlers) LoginUserHandler(c *gin.Context) {

	// Create an anonymous struct to hold the data that we expect to be in the request body
	var loginreq validator.LoginReq

	// Parse the request body into the anonymous struct
	err := utils.ReadRequest(c, &loginreq)
	if err != nil {
		utils.BadRequestResponse(c, "Authentication failed", 401, err)
		return
	}

	//perform validation for the loginrequest struct
	errorResponse, err := loginreq.Validate()

	if errorResponse != nil || err != nil {
		utils.ValidationErrorResponse(c, json.RawMessage(errorResponse))
		return
	}

	//retrieve user from the database
	user, err := h.App.Models.Users.GetByEmail(loginreq.Email)
	if err != nil {
		utils.LogError(c, err)
		utils.BadRequestResponse(c, "Authentication failed", 401, err)
		return
	}

	//check if password is correct
	if match := utils.CheckPasswordHash(loginreq.Password, user.Password); !match {
		utils.BadRequestResponse(c, "Authentication failed", 401, err)
		return
	}

	//generate jwt token
	accessToken, err := utils.GenerateJWT(user)
	if err != nil {
		utils.LogError(c, err)
		// app.errorResponse(c, http.StatusInternalServerError, envelope{"error": "Failed to generate JWT"})
		utils.BadRequestResponse(c, "Authentication failed", 401, err)
		return
	}

	//Respond with a 200 OK status code
	response := utils.Envelope{
		"status":  "success",
		"message": "Login successful",
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
	utils.WriteResponse(c, 200, response, nil)
}
