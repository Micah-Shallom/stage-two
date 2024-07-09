package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/Micah-Shallom/stage-two/utils"
	"github.com/Micah-Shallom/stage-two/validator"
	"github.com/gin-gonic/gin"
)

// addUserToOrganisationHandler is a handler function to add a user to an organization
func (h *Handlers) AddUserToOrganisationHandler(c *gin.Context) {
	orgID := c.Param("orgId")

	// Create an anonymous struct to hold the data that we expect to be in the request body
	var oauReq validator.OrgAddUserReq

	// Parse the request body into the anonymous struct
	err := utils.ReadRequest(c, &oauReq)
	if err != nil {
		utils.BadRequestResponse(c, "Client error", 400, err)
		return
	}

	//perform validation for the request struct
	errorResponse, err := oauReq.Validate()
	if errorResponse != nil || err != nil {
		utils.ValidationErrorResponse(c, json.RawMessage(errorResponse))
		return
	}

	//retrieve the authenticated user ID from context
	authUserID, exists := c.Get("UserID")
	if !exists {
		utils.BadRequestResponse(c, "Client error", 400, nil)
		return
	}

	//retrieve organization from the database
	org, err := h.App.Models.Organisations.GetByOrgID(orgID)
	if err != nil {
		utils.BadRequestResponse(c, "Client error", 400, err)
		return
	}

	//extract OrgID from the organization
	orgID = string(org.OrgID)

	//check if the authenticated user is a member of the organization
	exists, err = h.App.Models.Organisations.IsUserInOrganisation(orgID, authUserID.(string))
	if err != nil {
		utils.BadRequestResponse(c, "Client error", 400, err)
		return
	}

	fmt.Println(exists)
	if !exists {
		utils.BadRequestResponse(c, "Client error", 400, err)
		return
	}

	//check if the user to be added to the organization exists
	user, err := h.App.Models.Users.GetByID(oauReq.UserId)
	if err != nil {
		utils.BadRequestResponse(c, "Client error", 400, err)
		return
	}

	//add the user to the organization
	err = h.App.Models.Organisations.AddUserToOrganisation(orgID, user.UserID)
	if err != nil {
		utils.BadRequestResponse(c, "Client error", 400, err)
		return
	}

	//Respond with a 200 Created status code
	response := utils.Envelope{
		"status":  "success",
		"message": "User added to organisation successfully",
	}

	utils.WriteResponse(c, 200, response, nil)
}
