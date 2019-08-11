package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterAPI register web api
func RegisterAPI(router *mux.Router) {
	router.HandleFunc("/", Index)
	router.HandleFunc("/hello/:name", Hello)
}

// Index api index
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")

}

// Hello api hello
func Hello(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}
