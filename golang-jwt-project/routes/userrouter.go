package routes

import (
	controller "github.com/Preetham-P/go-projects/golang-jwt-project/controllers"
	middleware "github.com/Preetham-P/go-projects/golang-jwt-project/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.Use(middleware.Authenticate())
	router.GET("/users", controller.GetUsers())
	router.GET("/users/:user_id", controller.GetUserById())

}
