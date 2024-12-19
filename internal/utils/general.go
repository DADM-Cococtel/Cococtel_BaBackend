package utils

import (
	"github.com/Cococtel/Cococtel_BaBackend/internal/defines"
	"github.com/Cococtel/Cococtel_BaBackend/internal/domain/dtos"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/xlzd/gotp"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateHash(stringToHash string) (string, error) {
	stringHashed, err := bcrypt.GenerateFromPassword([]byte(stringToHash), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(stringHashed), nil
}

func CompareHashAndPassword(passwordHashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHashed), []byte(password))
}

func GenerateJWTToken(secret string, userID string, hours int) (string, error) {
	// Set custom claims
	duration := time.Duration(hours)
	claims := &dtos.JwtCustomClaims{
		User:   userID,
		Secret: secret,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * duration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var secretKey = []byte(os.Getenv(defines.JWTAuth))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return tokenString, nil
}

func GenerateTOTPWithSecret(username string, secret string) string {
	return gotp.NewDefaultTOTP(secret).ProvisioningUri(username, defines.AppName)
}

func GetUserIDFromToken(ctx *gin.Context) (string, error) {
	tokenString := ctx.GetHeader("x-auth-token")
	tokenString = tokenString[len("Bearer "):]
	token, _ := jwt.ParseWithClaims(tokenString, &dtos.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(defines.JWTAuth)), nil
	})
	if claims, ok := token.Claims.(*dtos.JwtCustomClaims); ok && token.Valid {
		return claims.User, nil
	}
	return "", defines.ErrInvalidToken
}
