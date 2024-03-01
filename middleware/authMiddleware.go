package middleware

import (
	"fmt"
	"net/http"
	"os"
	"task-5-pbi-btpns/database"
	"task-5-pbi-btpns/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *gin.Context) {
	// get the cookie
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Parse takes the token string and a function for looking up the key. The latter is especially
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		//check expired
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//find user with token
		var user models.UserModel
		if err := database.DB.First(&user, claims["user_id"]).Error; err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//attach to req
		c.Set("userId", claims["user_id"])

		//continue
		c.Next()

		// fmt.Println(claims["foo"], claims["exp"])
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
