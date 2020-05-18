package quic

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
	Network          string            `yarm:"network"`
	Token            string            `yarm:"token"`
	UID              uint64            `yarm:"uid"`
	ReadBufferSize   int               `yarm:"read_buffer_size"`
	WriteBufferSize  int               `yarm:"write_buffer_size"`
	ReadDeadline     time.Duration     `yarm:"read_dead_line"`
}

//ServerConfig grpc server config
type ServerConfig struct {
	Address         string        `yaml:"address"`
	MaxConnNum      int           `yaml:"max_conn_num"`
	SocketQueueSize int           `yaml:"socket_queue_size"`
	MaxMessageSize  uint32        `yaml:"max_message_size"`
	ReadBufferSize  int           `yaml:"read_buffer_size"`
	WriteBufferSize int           `yaml:"write_buffer_size"`
	ReadDeadline    time.Duration `yaml:"read_dead_line"`
	Network         string        `yaml:"network"`
}
