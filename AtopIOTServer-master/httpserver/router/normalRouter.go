package router

import (
	"errors"
	"net/http"

	"github.com/austinjan/AtopIOTServer/mongodb"
	"github.com/austinjan/AtopIOTServer/utils"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getVersion(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(r, w)
	res.AddPayload(bson.M{"ver": "1.0.0"})
	db := mongodb.GetDB()
	db.Test()
	res.SendResponse()

}

func formtest(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(r, w)
	body, err := parseBody(r)
	if err != nil {
		res.SendBodyError(err)
		return
	}
	res.AddPayload(body)
	res.SendResponse()

}

func getToken(w http.ResponseWriter, r *http.Request) {
	expireTime := viper.GetInt("expire")
	refreshExpireTime := viper.GetInt("refreshexpire")
	res := NewResponse(r, w)
	body, err := parseBody(r)
	if err != nil {
		res.SendBodyError(err)
		return
	}
	db := mongodb.GetDB()
	if err = mongodb.CheckKeysExist(body, []string{"name", "password"}); err != nil {
		res.SendBodyError(err)
		return
	}

	userid, err := db.ValidUser(body)
	if err != nil {
		res.SendGeneralError(errors.New("Account or password wrong"))
		return
	}
	token, err := utils.GetJWTToken(expireTime, userid)
	if err != nil {
		res.SendServerError(err)
		return
	}
	refreshToken, err := utils.GetJWTToken(refreshExpireTime, userid)
	if err != nil {
		res.SendServerError(err)
		return
	}
	id, _ := primitive.ObjectIDFromHex(userid)
	user, err := db.FindOne("users", bson.M{"_id": id})
	if err != nil {
		res.SendServerError(err)
		return
	}
	res.AddPayload(bson.M{"token": bson.M{"token": token, "refreshToken": refreshToken},
		"user": bson.M{"id": user["_id"], "name": user["name"]}})
	res.SendResponse()
}

func initNormalRouter(r *mux.Router) {
	// compatible for old api, not create
	r.HandleFunc("/api/version", getVersion).Methods("GET")
	r.HandleFunc("/api/formtest", formtest).Methods("POST")
	r.HandleFunc("/api/token", getToken).Methods("POST")

}
