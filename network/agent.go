package network

import (
	"net"
	"reflect"

	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/server"
)

//IAgent agent interface define
type IAgent interface {
	Run()
	OnClose()
	WriteMessage(msg interface{})
	ReadMessage() ([]byte, error)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
	Conn() Conn
	SetConn(conn Conn)
	Codec() codec.Codec
	SetCodec(c codec.Codec)
	Init(Conn, interface{}, codec.Codec)
	Token() string
	SetToken(string)
}

//NewAgent create new agent
func NewAgent(conn Conn, userData interface{}, codec codec.Codec) IAgent {
	return &Agent{conn: conn, userData: userData, codec: codec}
}

//Agent base agent struct
type Agent struct {
	conn     Conn
	userData interface{}
	codec    codec.Codec
	token    string
}

//Init agent
func (a *Agent) Init(conn Conn, userData interface{}, codec codec.Codec) {
	a.conn = conn
	a.userData = userData
	a.codec = codec
}

//Token get token
func (a *Agent) Token() string {
	return a.token
}

//SetToken set token
func (a *Agent) SetToken(token string) {
	a.token = token
}

//Conn get conn
func (a *Agent) Conn() Conn {
	return a.conn
}

// SetConn set conn
func (a *Agent) SetConn(conn Conn) {
	a.conn = conn
}

// Codec get codec
func (a *Agent) Codec() codec.Codec {
	return a.codec
}

//SetCodec set codec
func (a *Agent) SetCodec(c codec.Codec) {
	a.codec = c
}

//Run agent start run
func (a *Agent) Run() {
	// for {
	// 	data, err := a.ReadMessage()
	// 	if err != nil {
	// 		server.Debug("read message: %v", err)
	// 		break
	// 	}
	// 	out := make(map[string]interface{})
	// 	if err := a.Codec().Unmarshal(data, &out, nil); err == nil {
	// 		server.Info("receive message %v", out)
	// 		out["seq"] = out["seq"].(float64) + 1
	// 	} else {
	// 		server.Error(err)
	// 	}
	// 	a.WriteMessage(out)
	// }
}

func (a *Agent) OnClose() {

}

func (a *Agent) WriteMessage(msg interface{}) {
	data, err := a.Codec().Marshal(msg, nil)
	if err != nil {
		server.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
		return
	}
	err = a.conn.WriteMessage(data)
	if err != nil {
		server.Error("write message %v error: %v", reflect.TypeOf(msg), err)
	}

}

//ReadMessage goroutine not safe
func (a *Agent) ReadMessage() ([]byte, error) {
	b, err := a.conn.ReadMessage()
	return b, err
}

func (a *Agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *Agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *Agent) Close() {
	a.conn.Close()
}

func (a *Agent) Destroy() {
	a.conn.Destroy()
}

func (a *Agent) UserData() interface{} {
	return a.userData
}

func (a *Agent) SetUserData(data interface{}) {
	a.userData = data
}
