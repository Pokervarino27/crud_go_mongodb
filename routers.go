package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

//Route struct parametros del route
type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

//Routes arreglo del struct Route
type Routes []Route

//NewRouter funcion que retorna un router para conexion de un server
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandleFunc)
	}
	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"MovieList",
		"GET",
		"/movies",
		MovieList,
	},
	Route{
		"MovieAdd",
		"POST",
		"/movie",
		MovieAdd,
	},
	Route{
		"MovieShow",
		"GET",
		"/movie/{id}",
		MovieShow,
	},
	Route{
		"MovieUpdate",
		"PUT",
		"/movie/{id}",
		MovieUpdate,
	},
	Route{
		"MovieDelete",
		"DELETE",
		"/movie/{id}",
		MovieDelete,
	},
}
