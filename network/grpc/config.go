package grpc

import (
	"time"

	"github.com/gonethopper/nethopper/network/common"
)

//ClientConfig grpc client config
type ClientConfig struct {
	Nodes            []common.NodeInfo `mapstructure:"nodes"`
	ConnNum          int               `mapstructure:"conn_num"`
	ConnectInterval  time.Duration     `mapstructure:"connect_interval"`
	SocketQueueSize  int               `mapstructure:"socket_queue_size"`
	MaxMessageSize   uint32            `mapstructure:"max_message_size"`
	HandshakeTimeout time.Duration     `mapstructure:"handshake_timeout"`
	AutoReconnect    bool              `mapstructure:"auto_reconnect"`
}

//ServerConfig grpc server config
type ServerConfig struct {
	Address         string `mapstructure:"address"`
	MaxConnNum      int    `mapstructure:"max_conn_num"`
	SocketQueueSize int    `mapstructure:"socket_queue_size"`
	MaxMessageSize  uint32 `mapstructure:"max_message_size"`
}
