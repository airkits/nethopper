package tcp

import (
	"time"

	"github.com/airkits/nethopper/network/common"
)

//ClientConfig tcp client config
type ClientConfig struct {
	Nodes            []common.NodeInfo `mapstructure:"nodes"`
	ConnNum          int               `mapstructure:"conn_num"`
	ConnectInterval  time.Duration     `mapstructure:"connect_interval"`
	SocketQueueSize  int               `mapstructure:"socket_queue_size"`
	MaxMessageSize   uint32            `mapstructure:"max_message_size"`
	HandshakeTimeout time.Duration     `mapstructure:"handshake_timeout"`
	AutoReconnect    bool              `mapstructure:"auto_reconnect"`
	Network          string            `mapstructure:"network"`
	Token            string            `mapstructure:"token"`
	UID              uint64            `mapstructure:"uid"`
	ReadBufferSize   int               `mapstructure:"read_buffer_size"`
	WriteBufferSize  int               `mapstructure:"write_buffer_size"`
	ReadDeadline     time.Duration     `mapstructure:"read_dead_line"`
	QueueSize        int               `mapstructure:"queue_size"`
}

//GetQueueSize get module queue size
func (c *ClientConfig) GetQueueSize() int {
	return c.QueueSize
}

//ServerConfig tcp server config
type ServerConfig struct {
	Address         string        `mapstructure:"address"`
	Network         string        `mapstructure:"network"`
	MaxConnNum      int           `mapstructure:"max_conn_num"`
	SocketQueueSize int           `mapstructure:"socket_queue_size"`
	MaxMessageSize  uint32        `mapstructure:"max_message_size"`
	ReadBufferSize  int           `mapstructure:"read_buffer_size"`
	WriteBufferSize int           `mapstructure:"write_buffer_size"`
	ReadDeadline    time.Duration `mapstructure:"read_dead_line"`
	QueueSize       int           `mapstructure:"queue_size"`
}

//GetQueueSize get module queue size
func (s *ServerConfig) GetQueueSize() int {
	return s.QueueSize
}
