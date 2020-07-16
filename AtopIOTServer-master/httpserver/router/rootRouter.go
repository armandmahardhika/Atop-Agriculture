package router

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// type bson map[string]interface{}

func init() {
	log.Println("Init rootRouter")
}

// InitRouter  Init rootRouter
func InitRouter(r *mux.Router) {
	s := mux.NewRouter().PathPrefix("/apis").Subrouter().StrictSlash(true)

	n := negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(s))
	n.Use(negroni.NewRecovery())
	r.PathPrefix("/apis").Handler(n)
	initSecureRouter(s)
	initMqttRouter(s)
	initCrudRouter(s)
	initNormalRouter(r)
}
