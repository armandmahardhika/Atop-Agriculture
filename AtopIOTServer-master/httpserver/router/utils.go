package router // Get JSON body and transform to bson
import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

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
