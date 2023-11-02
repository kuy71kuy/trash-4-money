package utils

import (
	"github.com/golang-jwt/jwt"
	"os"
)

type MyClaims struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.StandardClaims
}

func ParseToken(tokenString string) *MyClaims {
	// Parse token
	token, _ := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims
	} else {
		claims.ID = 0
		claims.Name = ""
		claims.Role = ""
		return claims
	}
}
