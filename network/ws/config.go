package ws

import (
	"time"

	"github.com/gonethopper/nethopper/network/common"
)

//ClientConfig websocket client config
type ClientConfig struct {
	Nodes            []common.NodeInfo `mapstructure:"nodes"`
	ConnNum          int               `mapstructure:"conn_num"`
	ConnectInterval  time.Duration     `mapstructure:"connect_interval"`
	SocketQueueSize  int               `mapstructure:"socket_queue_size"`
	MaxMessageSize   uint32            `mapstructure:"max_message_size"`
	HandshakeTimeout time.Duration     `mapstructure:"handshake_timeout"`
	AutoReconnect    bool              `mapstructure:"auto_reconnect"`
	Token            string            `mapstructure:"token"`
	QueueSize        int               `mapstructure:"queue_size"`
}

//GetQueueSize get module queue size
func (c *ClientConfig) GetQueueSize() int {
	return c.QueueSize
}

//ServerConfig websocket server config
type ServerConfig struct {
	Address         string        `mapstructure:"address"`
	MaxConnNum      int           `mapstructure:"max_conn_num"`
	SocketQueueSize int           `mapstructure:"socket_queue_size"`
	MaxMessageSize  uint32        `mapstructure:"max_message_size"`
	HTTPTimeout     time.Duration `mapstructure:"http_timeout"`
	CertFile        string        `mapstructure:"cert_file"`
	KeyFile         string        `mapstructure:"key_file"`
	QueueSize       int           `mapstructure:"queue_size"`
}

//GetQueueSize get module queue size
func (s *ServerConfig) GetQueueSize() int {
	return s.QueueSize
}
