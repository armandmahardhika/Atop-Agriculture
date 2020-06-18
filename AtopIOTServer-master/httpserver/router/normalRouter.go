package router

import (
	"net/http"

	"github.com/austinjan/AtopIOTServer/mongodb"
	"github.com/gorilla/mux"
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
	//expireTime := viper.GetInt("expire")
	//refreshExpireTime := viper.GetInt("refreshexpire")
	res := NewResponse(r, w)
	body, err := parseBody(r)
	if err != nil {
		res.SendBodyError(err)
		return
	}
	db := mongodb.GetDB()
	if err = mongodb.CheckKeysExist(body, []string{"tracecode"}); err != nil {
		res.SendBodyError(err)
		//fmt.Println(body)
		return
	}
	userid, err := db.ValidUser(body)
	if err != nil {
		res.SendAuthenticateError("No Data Found")
		return
	}
	/**token, err := utils.GetJWTToken(expireTime, userid)
	if err != nil {
		res.SendServerError(err)
	}
	refreshToken, err := utils.GetJWTToken(refreshExpireTime, userid)
	if err != nil {
		res.SendServerError(err)
		"token": bson.M{"token": token, "refreshToken": refreshToken},

	}**/
	id, _ := primitive.ObjectIDFromHex(userid)
	data, err := db.FindOne("users", bson.M{"_id": id})
	if err != nil {
		res.SendServerError(err)
	}
	//add object to payload
	res.AddPayload(bson.M{"tracecode": data["tracecode"], "pit": data["pit"], "prt": data["prt"]})
	//"user": bson.M{"id": user["_id"], "name": user["name"], "data_1": user["data_1"], "data_2": user["data_2"]}})
	res.SendResponse()
}

func initNormalRouter(r *mux.Router) {
	// compatible for old api, not create
	r.HandleFunc("/api/version", getVersion).Methods("GET")
	r.HandleFunc("/api/formtest", formtest).Methods("POST")
	//r.HandleFunc("/api/token", getToken).Methods("POST")
	r.HandleFunc("/search/tracecode", getToken).Methods("POST")

}
