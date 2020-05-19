package kcp

import (
	"time"

	"github.com/gonethopper/nethopper/network/common"
)

//ClientConfig grpc client config
type ClientConfig struct {
	Nodes               []common.NodeInfo `mapstructure:"nodes"`
	ConnNum             int               `mapstructure:"conn_num"`
	ConnectInterval     time.Duration     `mapstructure:"connect_interval"`
	SocketQueueSize     int               `mapstructure:"socket_queue_size"`
	MaxMessageSize      uint32            `mapstructure:"max_message_size"`
	HandshakeTimeout    time.Duration     `mapstructure:"handshake_timeout"`
	AutoReconnect       bool              `mapstructure:"auto_reconnect"`
	Token               string            `mapstructure:"token"`
	UID                 uint64            `mapstructure:"uid"`
	UDPSocketBufferSize int               `mapstructure:"udp_socket_buffer_size"` //UDP listener socket buffer
	Dscp                int               `mapstructure:"dscp"`                   //set DSCP(6bit)
	Sndwnd              int               `mapstructure:"sndwnd"`                 //per connection UDP send window
	Rcvwnd              int               `mapstructure:"rcvwnd"`                 //per connection UDP recv window
	Mtu                 int               `mapstructure:"mtu"`                    //MTU of UDP packets, without IP(20) + UDP(8)
	Nodelay             int               `mapstructure:"nodelay"`                //ikcp_nodelay()
	Interval            int               `mapstructure:"interval"`               //ikcp_nodelay()
	Resend              int               `mapstructure:"resend"`                 //ikcp_nodelay()
	Nc                  int               `mapstructure:"nc"`                     //ikcp_nodelay()
	ReadDeadline        time.Duration     `mapstructure:"read_dead_line"`
}

//ServerConfig grpc server config
type ServerConfig struct {
	Address             string        `mapstructure:"address"`
	MaxConnNum          int           `mapstructure:"max_conn_num"`
	SocketQueueSize     int           `mapstructure:"socket_queue_size"`
	MaxMessageSize      uint32        `mapstructure:"max_message_size"`
	UDPSocketBufferSize int           `mapstructure:"udp_socket_buffer_size"` //UDP listener socket buffer
	ReadDeadline        time.Duration `mapstructure:"read_dead_line"`
	Dscp                int           `mapstructure:"dscp"`     //set DSCP(6bit)
	Sndwnd              int           `mapstructure:"sndwnd"`   //per connection UDP send window
	Rcvwnd              int           `mapstructure:"rcvwnd"`   //per connection UDP recv window
	Mtu                 int           `mapstructure:"mtu"`      //MTU of UDP packets, without IP(20) + UDP(8)
	Nodelay             int           `mapstructure:"nodelay"`  //ikcp_nodelay()
	Interval            int           `mapstructure:"interval"` //ikcp_nodelay()
	Resend              int           `mapstructure:"resend"`   //ikcp_nodelay()
	Nc                  int           `mapstructure:"nc"`       //ikcp_nodelay()
}
