package middleware

import (
	"github.com/Cococtel/Cococtel_BaBackend/internal/defines"
	"github.com/Cococtel/Cococtel_BaBackend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
)

func verifyToken(tokenString string) error {
	var secretKey = []byte(os.Getenv(defines.JWTAuth))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		log.Println(err)
		return err
	}

	if !token.Valid {
		return defines.ErrInvalidToken
	}

	return nil
}

func ProtectedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("x-auth-token")
		if tokenString == "" {
			utils.Error(c, http.StatusUnauthorized, defines.ErrNotFoundAuthHeader.Error())
			c.Abort()
			return
		}
		tokenString = tokenString[len("Bearer "):]
		err := verifyToken(tokenString)
		if err != nil {
			log.Println(err)
			utils.Error(c, http.StatusUnauthorized, defines.ErrInvalidToken.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}
