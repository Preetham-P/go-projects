package main

import (
	"os"

	"github.com/Preetham-P/go-projects/golang-jwt-project/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("port")

	if port == "" {
		port = "9090"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.Run(":" + port)

}
