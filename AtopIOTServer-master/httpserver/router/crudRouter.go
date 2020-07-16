package router

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/austinjan/AtopIOTServer/mongodb"
	"github.com/austinjan/AtopIOTServer/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

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
		res.SendAPIInputError(err)
		return
	}
	query, err := parseQuery(r)
	if err != nil {
		res.SendAPIInputError(err)
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
	res.AddPayload(bson.M{"data": data, "total": count})
	res.SendResponse()
}

// update handler
func putID(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(r, w)
	c, err := getURLCollection(r)
	if err != nil {
		res.SendAPIInputError(err)
		return
	}

	// Check permission
	rpi := requestToPermissionInfo(r, c)
	if pass := utils.CheckPermission(&rpi); !pass {
		res.SendPermissionError(rpi)
		return
	}

	fmt.Println("putID 2")

	id, err := getURLID(r)
	if err != nil {
		res.SendAPIInputError(err)
		return
	}

	body, err := parseBody(r)
	if err != nil {
		res.SendBodyError(err)
		return
	}
	delete(body, "_id")
	db := mongodb.GetDB()

	result, updateErr := db.UpdateID(c, id, body)
	if updateErr != nil {
		res.SendDatabaseError(updateErr)
		return
	}
	res.AddPayload(bson.M{"pre": result})
	res.SendResponse()

}

func post(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(r, w)
	c, err := getURLCollection(r)
	if err != nil {
		res.SendAPIInputError(err)
		return
	}

	rpi := requestToPermissionInfo(r, c)
	if pass := utils.CheckPermission(&rpi); !pass {
		res.SendPermissionError(rpi)
		return
	}

	body, err := parseBody(r)
	if err != nil {
		res.SendBodyError(err)
		return
	}
	delete(body, "_id")
	db := mongodb.GetDB()
	// check unique
	if pass, msg := db.CheckUniqueValue(c, body); !pass {
		res.SendUniqueValueError(errors.New(msg))
		return
	}

	result, err := db.Insert(c, body)
	if err != nil {
		res.SendDatabaseError(err)
		return
	}
	res.AddPayload(result)
	res.SendResponse()
}

func remove(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(r, w)
	c, err := getURLCollection(r)
	if err != nil {
		res.SendAPIInputError(err)
		return
	}

	rpi := requestToPermissionInfo(r, c)
	if pass := utils.CheckPermission(&rpi); !pass {
		res.SendPermissionError(rpi)
		return
	}

	id, err := getURLID(r)
	if err != nil {
		res.SendAPIInputError(err)
		return
	}

	db := mongodb.GetDB()

	result, err := db.DeleteOne(c, id)
	if err != nil {
		res.SendDatabaseError(err)
		return
	}
	res.AddPayload(result)
	res.SendResponse()
}

func deleteMany(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(r, w)
	c, err := getURLCollection(r)
	if err != nil {
		res.SendAPIInputError(err)
		return
	}

	query, err := parseQuery(r)
	if err != nil {
		res.SendAPIInputError(err)
		return
	}

	rpi := requestToPermissionInfo(r, c)
	if pass := utils.CheckPermission(&rpi); !pass {
		res.SendPermissionError(rpi)
	}

	db := mongodb.GetDB()
	result, err := db.DeleteMany(c, query)
	if err != nil {
		res.SendDatabaseError(err)
		return
	}
	res.AddPayload(result)
	res.SendResponse()
}

// params: {ids:[_id1, _id2]}  _id should be string
func deleteIDs(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(r, w)
	c, err := getURLCollection(r)
	if err != nil {
		res.SendAPIInputError(err)
		return
	}

	querys, err := url.ParseQuery(r.URL.RawQuery)

	query := make(map[string]interface{})
	for k, v := range querys {
		query[k] = v
	}

	rpi := requestToPermissionInfo(r, c)
	if pass := utils.CheckPermission(&rpi); !pass {
		res.SendPermissionError(rpi)
	}

	db := mongodb.GetDB()
	result, err := db.DeleteIDs(c, query)
	if err != nil {
		res.SendDatabaseError(err)
		return
	}
	res.AddPayload(result)
	res.SendResponse()
}

func initCrudRouter(r *mux.Router) {
	r.HandleFunc("/ping", ping).Methods("GET")
	r.HandleFunc("/crud/{collection}", read).Methods("GET")
	r.HandleFunc("/crud/{collection}/{id:[0-9a-f]+}", putID).Methods("PUT")
	r.HandleFunc("/crud/{collection}", post).Methods("POST")
	r.HandleFunc("/crud/{collection}/{id:[0-9a-f]+}", remove).Methods("DELETE")
	r.HandleFunc("/crud/{collection}", deleteMany).Methods("DELETE")
	r.HandleFunc("/delete/{collection}", deleteIDs).Methods("DELETE")
}
