package router

import (
	"fmt"
	"net/http"

	"github.com/austinjan/AtopIOTServer/mongodb"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getUserID(r *http.Request) string {
	user := r.Context().Value("user")
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	key, ok := claims["user"]
	if !ok {
		return ""
	}
	return key.(string)
}

func ping(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(r, w)
	res.AddPayload(bson.M{"message": "pong"})
	user := r.Context().Value("user")
	fmt.Fprintf(w, "This is an authenticated request")
	fmt.Fprintf(w, "Claim content:\n")
	for k, v := range user.(*jwt.Token).Claims.(jwt.MapClaims) {
		fmt.Fprintf(w, "%s :\t%#v\n", k, v)
	}
	res.SendResponse()

}

// get handler
func read(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(r, w)
	c, err := getURLCollection(r)
	if err != nil {
		res.SendGeneralError(err)
		return
	}
	query, err := parseQuery(r)
	if err != nil {
		res.SendQueryError(err)
		return
	}
	db := mongodb.GetDB()
	data, err := db.Read(c, query)
	if err != nil {
		res.SendServerError(err)
		return
	}
	count, err := db.Count(c)
	if err != nil {
		res.SendServerError(err)
		return
	}
	if count <= 0 {
		res.SendEmptyError()
		return
	}
	res.AddPayload(bson.M{"data": data, "count": count})
	res.SendResponse()
}

// update handler
func putID(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(r, w)
	c, err := getURLCollection(r)
	if err != nil {
		res.SendGeneralError(err)
		return
	}
	id, err := getURLID(r)
	if err != nil {
		res.SendGeneralError(err)
		return
	}
	body, err := parseBody(r)
	if err != nil {
		res.SendBodyError(err)
		return
	}
	db := mongodb.GetDB()
	result, updateErr := db.UpdateID(c, id, body)
	if updateErr != nil {
		res.SendGeneralError(updateErr)
		return
	}
	res.AddPayload(bson.M{"pre": result})
	res.SendResponse()

}

func getCurrentUser(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(r, w)
	idHex := getUserID(r)
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		res.SendServerError(err)
		return
	}
	db := mongodb.GetDB()
	result, err := db.FindOne("users", bson.M{"_id": id})
	if err != nil {
		res.SendServerError(err)
		return
	}
	res.AddPayload(bson.M{"user": result})
	res.SendResponse()
}

func initCrudRouter(r *mux.Router) {

	r.HandleFunc("/ping", ping).Methods("GET")
	r.HandleFunc("/crud/{collection}", read).Methods("GET")
	r.HandleFunc("/crud/{collection}/{id:[0-9a-f]+}", putID).Methods("PUT")
	r.HandleFunc("/currentUser", getCurrentUser).Methods("GET")
}
