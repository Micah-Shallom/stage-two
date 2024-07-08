package routes

import (
	"github.com/Micah-Shallom/stage-two/handlers"
	"github.com/Micah-Shallom/stage-two/middleware"
	"github.com/gin-gonic/gin"
)

func SetupOrganisationRoutes(r *gin.Engine, handler *handlers.Handlers) {
	r.GET("/api/organisations", middleware.JWTMiddleware(), handler.GetOrganizationsHandler)
	r.GET("/api/organisations/:orgId", middleware.JWTMiddleware(), handler.GetOrganisationByIDHandler)
	r.POST("/api/organisations", middleware.JWTMiddleware(), handler.CreateOrganisationHandler)
	r.POST("/api/organisations/:orgId/users", middleware.JWTMiddleware(), handler.AddUserToOrganisationHandler)
}
