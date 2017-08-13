package server

import (
	"fmt"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.
		Methods("GET").
		Path(fmt.Sprintf("/race/{%s}/{%s}", startVar, endVar)).
		HandlerFunc(RaceWithTitles)
	return router
}
