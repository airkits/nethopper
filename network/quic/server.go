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

	"github.com/airkits/nethopper/config"
	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/network"
	quic "github.com/lucas-clemente/quic-go"
)

//NewServer create quic server
func NewServer(conf config.IConfig, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) network.IServer {
	s := new(Server)
	s.Conf = conf.(*ServerConfig)
	s.NewAgent = agentFunc
	s.CloseAgent = agentCloseFunc
	s.wg = &sync.WaitGroup{}
	return s
}

// Server quic server define
type Server struct {
	Conf       *ServerConfig
	NewAgent   network.AgentCreateFunc
	listener   quic.Listener
	CloseAgent network.AgentCloseFunc
	conns      ConnSet
	mutexConns sync.Mutex
	wg         *sync.WaitGroup
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

	listener, err := quic.ListenAddr(s.Conf.Address, s.generateTLSConfig(), nil)
	if err != nil {
		panic(err)
	}
	s.listener = listener

	log.Info("listening on: %s %s", s.Conf.Network, listener.Addr())

	// loop accepting
	for {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			log.Warning("accept failed: %s", err.Error())
			continue
		}
		stream, err := sess.AcceptStream(context.Background())
		if err != nil {
			panic(err)

		}
		log.Info("receive one client peer %s", sess.RemoteAddr().String())
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
	c := NewConn(sess, stream, s.Conf.SocketQueueSize, s.Conf.MaxMessageSize, s.Conf.ReadDeadline*time.Second)
	agent = s.NewAgent(c, 0, "")
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
