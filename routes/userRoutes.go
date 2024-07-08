package routes

import (
	"github.com/Micah-Shallom/stage-two/handlers"
	"github.com/Micah-Shallom/stage-two/middleware"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, handler *handlers.Handlers) {
	r.GET("/api/users/:id",middleware.JWTMiddleware(), handler.GetUserHandler)
}
