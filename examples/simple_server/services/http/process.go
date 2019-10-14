package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/examples/simple_server/common"
	"github.com/gonethopper/nethopper/examples/simple_server/pb"
	"github.com/gonethopper/nethopper/server"
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

	sbody, e := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if e != nil {
		w.WriteHeader(500)
		return
	}
	server.Info(string(sbody))
	var v = make(map[string]interface{})
	if err := codec.JSONCodec.Unmarshal(sbody, &v, nil); err != nil {
		server.Info(err)
		return
	}
	body := &pb.User{
		Uid:    string(int64(v["uid"].(float64))),
		Passwd: "",
	}
	sess := server.GetSession(token)
	if sess != nil {
		m := server.CreateMessage(common.MessageIDLogin, server.ServiceIDHTTP, server.ServiceIDLogic, server.MTRequest, common.MessageIDLoginCmd, token)
		m.SetBody(body)
		server.SendMessage(m.DestID, 0, m)
	}
	defer close(sess.Die)

	result := <-sess.Done //等待Done的通知，此时call.Reply发生了变化。

	respBody := (result.Response).(*pb.User)
	server.Info("message done,get pwd  %s", respBody.String())
	fmt.Fprint(w, string(respBody.Passwd))

	// var i int
	// for start := time.Now(); ; {

	// 	if i>>3 == 1 {
	// 		i = 1
	// 		if time.Since(start) > TIMEOUT {
	// 			fmt.Fprint(w, "timeout")
	// 			return
	// 		}
	// 		runtime.Gosched()
	// 	}
	// 	i++

	// 	if v, err := sess.MQ.AsyncPop(); err == nil {
	// 		fmt.Fprint(w, v.(server.Message).Payload)
	// 		fmt.Fprint(w, "close bybye")

	// 		return
	// 	} else if err == queue.ErrQueueIsClosed {
	// 		fmt.Fprint(w, err.Error())
	// 		return
	// 	}
	// }

}

// Hello api hello
func Hello(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}
