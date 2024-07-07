package handlers

import (
	"github.com/Micah-Shallom/stage-two/utils"
	"github.com/gin-gonic/gin"
)

// getOrganizationsHandler is a handler function to get organizations
func (h *Handlers) GetOrganizationsHandler(c *gin.Context) {
	//retrieve the authenticated user ID from context
	authUserID, exists := c.Get("UserID")
	if !exists {
		utils.BadRequestResponse(c, "Client error", 400, nil)
		return
	}

	//Fetch organizations the authenticted user belongs to or has created
	orgs, err := h.App.Models.Organisations.GetByUserID(authUserID.(string))
	if err != nil {
		utils.BadRequestResponse(c, "Client error", 400, err)
		return
	}

	// prepare response payload
	var orgResponse []map[string]any
	for _, org := range orgs {
		orgResponse = append(orgResponse, map[string]any{
			"orgId":       org.OrgID,
			"name":        org.Name,
			"description": org.Description,
		})
	}
	utils.SendUserOrganisationsResponse(c, orgResponse)
}
