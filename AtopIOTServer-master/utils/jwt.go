package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)
// SecretKey secretkey used in jwt signed
var SecretKey = "atopsecretkey"

// GetJWTToken get token
func GetJWTToken(expireTime int, id string) (string, error) {
	//	expireTime := 3600
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Second * time.Duration(expireTime)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["user"] = id
	claims["ping"] = "pong"
	token.Claims = claims
	jwtString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return jwtString, nil
}
