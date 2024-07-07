package handlers

import (
	"github.com/Micah-Shallom/stage-two/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetOrganisationByIDHandler(c *gin.Context) {
	orgID := c.Param("orgId")
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
	if !exists {
		utils.BadRequestResponse(c, "Client error", 400, err)
		return
	}

	utils.SendOrganisationResponse(c, org)
}
