package handlers

import (
	"net/http"

	"github.com/Micah-Shallom/stage-two/models"
	"github.com/Micah-Shallom/stage-two/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) CreateOrganisationHandler(c *gin.Context) {
	// Create an anonymous struct to hold the data that we expect to be in the request body
	var input struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
	}

	// Parse the request body into the anonymous struct
	valid, errors := utils.ReadRequest(c, &input)
	if !valid {
		utils.ValidationErrorResponse(c, errors)
		return
	}

	//retrieve the authenticated user ID from context
	authUserID, exists := c.Get("UserID")
	if !exists {
		utils.BadRequestResponse(c, "Client error", 400, nil)
		return
	}

	// Create a new Organisation object containing the data that we read from the request body
	org := &models.Organisation{
		OrgID:       utils.GenerateUUID(),
		Name:        input.Name,
		Description: input.Description,
	}

	//save organisation model to the database
	err := h.App.Models.Organisations.Create(org)
	if err != nil {
		utils.BadRequestResponse(c, "Client error", 400, err)
		return
	}

	//add the authenticated user to the organisation
	err = h.App.Models.Organisations.AddUserToOrganisation(org.OrgID, authUserID.(string))
	if err != nil {
		utils.BadRequestResponse(c, "Client error", 400, err)
		return
	}

	//Respond with a 201 Created status code
	response := utils.Envelope{
		"status":  "success",
		"message": "Organisation created successfully",
		"data": map[string]any{
			"orgId":       org.OrgID,
			"name":        org.Name,
			"description": org.Description,
		},
	}
	utils.WriteResponse(c, http.StatusCreated, response, nil)

}
