package routes

import (
	"github.com/Micah-Shallom/stage-two/handlers"
	"github.com/Micah-Shallom/stage-two/middleware"
	"github.com/gin-gonic/gin"
)

// routes sets up the application routes
func Routes(r *gin.Engine, handler *handlers.Handlers) *gin.Engine {

	//public routes
	SetupAuthRoutes(r, handler)
	//custom error handlers
	SetupErrorRoutes(r)

	//middleware to protect routes
	r.Use(middleware.JWTMiddleware())

	//protected routes
	//USER routes
	SetupUserRoutes(r, handler)

	//ORGANISATION routes
	SetupOrganisationRoutes(r, handler)

	return r
}
