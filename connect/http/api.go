package http

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// RegisterAPI register web api
func RegisterAPI(router *httprouter.Router) {
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
}

// Index api index
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")

}

// Hello api hello
func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}
