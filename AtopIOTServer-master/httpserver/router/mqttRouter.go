package router

import "github.com/gorilla/mux"

func initMqttRouter(r *mux.Router) {
	r.HandleFunc("/mqtt/subscribes", getSubscribes).Methods("GET")

}
