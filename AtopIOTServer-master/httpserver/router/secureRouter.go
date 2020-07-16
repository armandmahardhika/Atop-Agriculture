package router

import (
	"fmt"
	"net/http"

	"github.com/austinjan/AtopIOTServer/cmd"
	"github.com/austinjan/AtopIOTServer/mongodb"
	"github.com/austinjan/AtopIOTServer/utils"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getCurrentUser(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(r, w)
	idHex := GetUserID(r)
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

// RefreshResponse : Responese json for /apis/v1/refresh
type RefreshResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

// api rul format : /apis/[version]/refresh
// return body: array of json for data
func refreshToken(w http.ResponseWriter, r *http.Request) {

	expireTime := viper.GetInt("expire")
	refreshExpireTime := viper.GetInt("refreshexpire")
	fmt.Println("Refresh token", expireTime, refreshExpireTime)
	userid := GetUserID(r)
	res := NewResponse(r, w)

	token, err := utils.GetJWTToken(expireTime, userid)
	if err != nil {
		res.SendServerError(err)
	}
	refreshToken, err := utils.GetJWTToken(refreshExpireTime, userid)
	if err != nil {
		res.SendServerError(err)
	}

	res.AddPayload(bson.M{"token": bson.M{"token": token, "refreshToken": refreshToken}})
	res.SendResponse()
}

func getCollections(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(r, w)
	db := mongodb.GetDB()
	collections := db.GetCollections()
	res.AddPayload(bson.M{"collections": collections})
	res.SendResponse()
}

func fuzzySearchDocument(w http.ResponseWriter, r *http.Request) {
	//_, collection, err := bodyParser(r)
	res := NewResponse(r, w)

	collection, err := getURLCollection(r)
	if err != nil {
		res.SendGeneralError(err)
		return
	}
	searchText, err := getURLParam(r, "string")
	if err != nil {
		res.SendGeneralError(err)
		return
	}

	db := mongodb.GetDB()
	result, err := db.FuzzySearch(collection, searchText)
	res.AddPayload(bson.M{"result": result})
	res.SendResponse()
}

func mqttBrokerInfo(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(r, w)
	info, err := cmd.GetMosquittoInfo()
	if err != nil {
		res.SendGeneralError(err)
		return
	}
	res.AddPayload(bson.M{"info": info})
	res.SendResponse()
}

func initSecureRouter(r *mux.Router) {
	r.HandleFunc("/currentUser", getCurrentUser).Methods("GET")
	r.HandleFunc("/refresh/token", refreshToken).Methods("GET")
	r.HandleFunc("/collections", getCollections).Methods("GET")
	r.HandleFunc("/search/{collection}/{string}", fuzzySearchDocument).Methods("GET")
	r.HandleFunc("/mqtt/broker/info", mqttBrokerInfo).Methods("GET")
}
