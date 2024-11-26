package middleware

import (
	"fmt"
	"net/http"

	"github.com/Preetham-P/go-projects/golang-jwt-project/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientToken := ctx.Request.Header.Get("token")
		if clientToken == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("No Auth Header"),
			})
			return
		}
		claims, err := helpers.ValidateToken(clientToken)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		ctx.Set("email", claims.Email)
		ctx.Set("first_name", claims.First_Name)
		ctx.Set("last_name", claims.Last_Name)
		ctx.Set("user_id", claims.User_Id)
		ctx.Set("user_type", claims.User_Type)
		ctx.Next()

	}

}
