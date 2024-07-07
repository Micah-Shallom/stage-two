package utils

import (
	"net/http"

	"github.com/Micah-Shallom/stage-two/models"
	"github.com/gin-gonic/gin"
)

// sendUserResponse sends a response containing user data
func SendUserResponse(c *gin.Context, user *models.User) {
	response := Envelope{
		"status":  "success",
		"message": "User retrieved successfully",
		"data": map[string]any{
			"userId":    user.UserID,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
			"phone":     user.Phone,
		},
	}
	WriteResponse(c, http.StatusOK, response, nil)
}

// sendUserOrganisationsResponse sends a response containing user organisations
func SendUserOrganisationsResponse(c *gin.Context, orgResponse []map[string]any) {
	response := Envelope{
		"status":  "success",
		"message": "User organisations retrieved successfully",
		"data": map[string]any{
			"organisations": orgResponse,
		},
	}
	WriteResponse(c, http.StatusOK, response, nil)
}

// sendUserOrganisationResponse sends a response containing user organisation data
func SendOrganisationResponse(c *gin.Context, org *models.Organisation) {
	response := Envelope{
		"status":  "success",
		"message": "Organisation retrieved successfully",
		"data": map[string]any{
			"orgId":       org.OrgID,
			"name":        org.Name,
			"description": org.Description,
		},
	}
	WriteResponse(c, http.StatusOK, response, nil)
}
