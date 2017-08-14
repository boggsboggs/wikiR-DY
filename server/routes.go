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
	router.
		Methods("GET").
		Path(fmt.Sprintf("/race/url/{%s}/{%s}", startVar, endVar)).
		HandlerFunc(RaceWithURLs)
	return router
}
