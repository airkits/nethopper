package kcp

import (
	"time"

	"github.com/gonethopper/nethopper/network/common"
)

//ClientConfig grpc client config
type ClientConfig struct {
	Nodes               []common.NodeInfo `yarm:"nodes"`
	ConnNum             int               `yarm:"conn_num"`
	ConnectInterval     time.Duration     `yarm:"connect_interval"`
	SocketQueueSize     int               `yarm:"socket_queue_size"`
	MaxMessageSize      uint32            `yarm:"max_message_size"`
	HandshakeTimeout    time.Duration     `yarm:"handshake_timeout"`
	AutoReconnect       bool              `yarm:"auto_reconnect"`
	Token               string            `yarm:"token"`
	UID                 uint64            `yarm:"uid"`
	UDPSocketBufferSize int               `yarm:"udp_socket_buffer_size"` //UDP listener socket buffer
	Dscp                int               `yarm:"dscp"`                   //set DSCP(6bit)
	Sndwnd              int               `yarm:"sndwnd"`                 //per connection UDP send window
	Rcvwnd              int               `yarm:"rcvwnd"`                 //per connection UDP recv window
	Mtu                 int               `yarm:"mtu"`                    //MTU of UDP packets, without IP(20) + UDP(8)
	Nodelay             int               `yarm:"nodelay"`                //ikcp_nodelay()
	Interval            int               `yarm:"interval"`               //ikcp_nodelay()
	Resend              int               `yarm:"resend"`                 //ikcp_nodelay()
	Nc                  int               `yarm:"nc"`                     //ikcp_nodelay()
	ReadDeadline        time.Duration     `yarm:"read_dead_line"`
}

//ServerConfig grpc server config
type ServerConfig struct {
	Address             string        `yaml:"address"`
	MaxConnNum          int           `yaml:"max_conn_num"`
	SocketQueueSize     int           `yaml:"socket_queue_size"`
	MaxMessageSize      uint32        `yaml:"max_message_size"`
	UDPSocketBufferSize int           `yaml:"udp_socket_buffer_size"` //UDP listener socket buffer
	ReadDeadline        time.Duration `yaml:"read_dead_line"`
	Dscp                int           `yaml:"dscp"`     //set DSCP(6bit)
	Sndwnd              int           `yaml:"sndwnd"`   //per connection UDP send window
	Rcvwnd              int           `yaml:"rcvwnd"`   //per connection UDP recv window
	Mtu                 int           `yaml:"mtu"`      //MTU of UDP packets, without IP(20) + UDP(8)
	Nodelay             int           `yaml:"nodelay"`  //ikcp_nodelay()
	Interval            int           `yaml:"interval"` //ikcp_nodelay()
	Resend              int           `yaml:"resend"`   //ikcp_nodelay()
	Nc                  int           `yaml:"nc"`       //ikcp_nodelay()
}
