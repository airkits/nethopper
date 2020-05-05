package quic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"sync"
	"time"

	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
	quic "github.com/lucas-clemente/quic-go"
)

//NewServer create quic server
func NewServer(m map[string]interface{}, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) network.IServer {
	s := new(Server)
	if err := s.ReadConfig(m); err != nil {
		panic(err)
	}
	s.NewAgent = agentFunc
	s.CloseAgent = agentCloseFunc

	return s
}

// Server quic server define
type Server struct {
	Config
	NewAgent   network.AgentCreateFunc
	listener   quic.Listener
	CloseAgent network.AgentCloseFunc
	conns      ConnSet
	mutexConns sync.Mutex
	wg         sync.WaitGroup
}

// ReadConfig config map
// m := map[string]interface{}{
// readBufferSize default 32767
// writeBufferSize default 32767
// address default :16000
// network default "quic4"  use "quic4/quic6"
// readDeadline default 15
//	"maxConnNum":1024,
//  "socketQueueSize":100,
//  "maxMessageSize":4096
// //tls support
//  "certFile":"",
//  "keyFile":"",
// }
func (s *Server) ReadConfig(m map[string]interface{}) error {

	if err := server.ParseConfigValue(m, "address", ":16000", &s.Address); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "maxConnNum", 1024, &s.MaxConnNum); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "socketQueueSize", 100, &s.RWQueueSize); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "maxMessageSize", 4096, &s.MaxMessageSize); err != nil {
		return err
	}

	if err := server.ParseConfigValue(m, "readBufferSize", 32767, &s.ReadBufferSize); err != nil {
		return err
	}

	if err := server.ParseConfigValue(m, "writeBufferSize", 32767, &s.WriteBufferSize); err != nil {
		return err
	}

	if err := server.ParseConfigValue(m, "network", "quic4", &s.Network); err != nil {
		return err
	}

	if err := server.ParseConfigValue(m, "readDeadline", 15, &s.ReadDeadline); err != nil {
		return err
	}
	s.ReadDeadline = s.ReadDeadline * time.Second
	return nil
}

// Setup a bare-bones TLS config for the server
func (s *Server) generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}

//ListenAndServe start serve
func (s *Server) ListenAndServe() {

	listener, err := quic.ListenAddr(s.Address, s.generateTLSConfig(), nil)
	if err != nil {
		panic(err)
	}
	s.listener = listener

	server.Info("listening on: %s %s", s.Network, listener.Addr())

	// loop accepting
	for {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			server.Warning("accept failed: %s", err.Error())
			continue
		}
		stream, err := sess.AcceptStream(context.Background())
		if err != nil {
			panic(err)

		}
		server.Info("receive one client peer %s", sess.RemoteAddr().String())
		//go conn
		go s.Transport(sess, stream)
	}
}

//Transport quic connection
func (s *Server) Transport(sess quic.Session, stream quic.Stream) error {

	s.wg.Add(1)
	defer s.wg.Done()

	// s.conns[stream] = struct{}{}
	// s.mutexConns.Unlock()

	var agent network.IAgent
	c := NewConn(sess, stream, s.RWQueueSize, s.MaxMessageSize, s.ReadDeadline)
	agent = s.NewAgent(c, "")
	agent.Run()

	// cleanup
	stream.Context().Done()
	// s.mutexConns.Lock()
	// delete(s.conns, stream)
	// s.mutexConns.Unlock()
	s.CloseAgent(agent)
	agent.OnClose()
	return nil
}

//Close quic server
func (s *Server) Close() {
	s.listener.Close()

	s.mutexConns.Lock()
	for conn := range s.conns {
		conn.Context().Done()
	}
	s.conns = nil
	s.mutexConns.Unlock()

	s.wg.Wait()
}
