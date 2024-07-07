package routes

import (
	"github.com/Micah-Shallom/stage-two/utils"
	"github.com/gin-gonic/gin"
)

func SetupErrorRoutes(r *gin.Engine) {
	r.NoRoute(utils.NotFoundResponse)
	r.NoMethod(utils.MethodNotAllowedResponse)
}
