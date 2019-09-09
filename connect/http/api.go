package http

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gonethopper/nethopper/server"
	"github.com/gonethopper/queue"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

const TIMEOUT = time.Second * 15

// RegisterAPI register web api
func RegisterAPI(router *mux.Router) {
	router.HandleFunc("/", Index)
	router.HandleFunc("/hello/:name", Hello)
}

// Index api index
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
	token := context.Get(r, "token").(string)
	fmt.Fprint(w, token)
	sess := server.GetSession(token)
	if sess != nil {
		m := server.CreateMessage(server.ServiceIDHTTP, server.ServiceIDLogic, server.MTRequest, "index", []byte{})
		server.SendMessage(m.DestID, 0, m)
	}
	defer close(sess.Die)
	var i int
	for start := time.Now(); ; {

		if i>>3 == 1 {
			i = 1
			if time.Since(start) > TIMEOUT {
				fmt.Fprint(w, "timeout")
				return
			}
			runtime.Gosched()
		}
		i++

		if v, err := sess.MQ.AsyncPop(); err == nil {
			fmt.Fprint(w, v.(server.Message).Payload)
			fmt.Fprint(w, "close bybye")

			return
		} else if err == queue.ErrQueueIsClosed {
			fmt.Fprint(w, err.Error())
			return
		}
	}

}

// Hello api hello
func Hello(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}
