package router

import (
	"net/http"

	"github.com/austinjan/AtopIOTServer/utils"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.SecretKey), nil
	},
	ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
		res := NewResponse(r, w)
		res.SendAuthenticateError(err)
	},
	SigningMethod: jwt.SigningMethodHS256,
})
