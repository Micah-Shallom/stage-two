package routes

import (
	"github.com/Micah-Shallom/stage-two/handlers"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine, handler *handlers.Handlers) {
	r.POST("/auth/register", handler.RegisterUserHandler)
	r.POST("/auth/login", handler.LoginUserHandler)
}
