package router // Get JSON body and transform to bson
import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/austinjan/AtopIOTServer/mongodb"
	"github.com/austinjan/AtopIOTServer/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

//GetUserID  get user id from request jwt token
func GetUserID(r *http.Request) string {
	user := r.Context().Value("user")
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	key, ok := claims["user"]
	if !ok {
		return ""
	}
	return key.(string)
}

func parseBody(r *http.Request) (bson.M, error) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		return nil, err
	}

	var bodyBson map[string]interface{}

	if err := json.Unmarshal(b, &bodyBson); err != nil {
		return nil, err
	}

	return bodyBson, err
}

// ErrorResponse error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// parseQuery, used in normal query
func parseQuery(r *http.Request) (bson.M, error) {
	querys, err := url.ParseQuery(r.URL.RawQuery)

	query := make(map[string]interface{})
	for k, v := range querys {
		if len(v) == 1 {
			query[k] = v[0]
		} else {
			query[k] = v
		}
	}
	return query, err
}

// getURLCOllection(r,name) get url param by name
// example:
// collection, err := getURLParam(r,"ids")
// 	if err != nil {
// 		res.SendGeneralError(err)
// 		return
// 	}
func getURLParam(r *http.Request, param string) (string, error) {
	vars := mux.Vars(r)

	c, ok := vars[param]
	if !ok {
		return "", fmt.Errorf("Can not find {%s} in URL", param)
	}

	return c, nil
}

// getURLCOllection() get url collection template
// example:
// collection, err := getURLCollection(r)
// 	if err != nil {
// 		res.SendGeneralError(err)
// 		return
// 	}
func getURLCollection(r *http.Request) (string, error) {
	vars := mux.Vars(r)

	c, ok := vars["collection"]
	if !ok {
		return "", errors.New("URL format error")
	}

	return c, nil
}

func getURLID(r *http.Request) (string, error) {
	vars := mux.Vars(r)

	c, ok := vars["id"]
	if !ok {
		return "", errors.New("URL format error")
	}

	return c, nil
}

func requestToPermissionInfo(r *http.Request, c string) utils.RequestPermission {
	db := mongodb.GetDB()
	role, err := db.GetRoleByID(GetUserID(r))
	if err != nil {
		role.Role = ""
		role.Name = ""
	}
	permissinInfo := utils.RequestPermission{
		Role:       role.Role,
		Name:       role.Name,
		Method:     r.Method,
		Collection: c,
	}
	return permissinInfo
}
