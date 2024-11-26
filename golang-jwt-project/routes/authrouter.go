package routes

import (
	controller "github.com/Preetham-P/go-projects/golang-jwt-project/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {

	router.POST("/auth/signup", controller.SignUp())
	router.POST("/auth/login", controller.Login())

}
