package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ApiRoute struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

var routes = []ApiRoute{
	ApiRoute{
		"GET",
		"/",
		GetIndex,
	},
	ApiRoute{
		"PUT",
		"/api/{name}/",
		PutRecord,
	},
	ApiRoute{
		"POST",
		"/api/{name}/",
		PostRecord,
	},
	ApiRoute{
		"GET",
		"/api/",
		GetAllRecords,
	},
	ApiRoute{
		"GET",
		"/api/{name}/",
		GetRecord,
	},
	ApiRoute{
		"DELETE",
		"/api/",
		DeleteAllRecords,
	},
	ApiRoute{
		"DELETE",
		"/api/{name}/",
		DeleteRecord,
	},
}

func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		log.Printf(
			"%s  %s  %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler = route.Handler
		handler = Logger(handler)
		router.
			Methods(route.Method).
			Path(route.Path).
			Handler(handler)
	}
	return router
}
