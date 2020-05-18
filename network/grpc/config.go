package grpc

import (
	"time"

	"github.com/gonethopper/nethopper/network/common"
)

//ClientConfig grpc client config
type ClientConfig struct {
	Nodes            []common.NodeInfo `yarm:"nodes"`
	ConnNum          int               `yarm:"conn_num"`
	ConnectInterval  time.Duration     `yarm:"connect_interval"`
	SocketQueueSize  int               `yarm:"socket_queue_size"`
	MaxMessageSize   uint32            `yarm:"max_message_size"`
	HandshakeTimeout time.Duration     `yarm:"handshake_timeout"`
	AutoReconnect    bool              `yarm:"auto_reconnect"`
}

//ServerConfig grpc server config
type ServerConfig struct {
	Address         string `yaml:"address"`
	MaxConnNum      int    `yaml:"max_conn_num"`
	SocketQueueSize int    `yaml:"socket_queue_size"`
	MaxMessageSize  uint32 `yaml:"max_message_size"`
}
