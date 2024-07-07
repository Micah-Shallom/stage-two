package routes

import (
	"github.com/Micah-Shallom/stage-two/handlers"
	"github.com/gin-gonic/gin"
)

func SetupOrganisationRoutes(r *gin.Engine, handler *handlers.Handlers) {
	r.GET("/api/organisations", handler.GetOrganizationsHandler)
	r.GET("/api/organisations/:orgId", handler.GetOrganisationByIDHandler)
	r.POST("/api/organisations", handler.CreateOrganisationHandler)
	r.POST("/api/organisations/:orgId/users", handler.AddUserToOrganisationHandler)
}
