package mediator

const (
	// ModuleNamedID module id define, system reserved 1-63
	ModuleNamedID = iota
	// Main main goruntinue
	Main
	// Monitor server monitor module
	Monitor
	// Log log module
	Log
	// TCP tcp module
	TCP
	// KCP kcp module
	KCP
	// QUIC quic module
	QUIC
	// WSServer ws server
	WSServer
	// GRPCServer grpc server
	GRPCServer
	// HTTP http module
	HTTP
	// Logic logic module
	Logic
	// Redis redis module
	Redis
	// TCPClient tcp client module
	TCPClient
	// KCPClient kcp client module
	KCPClient
	// QUICClient quic client module
	QUICClient
	// HTTPClient http client module
	HTTPClient
	// GRPCClient grpc client module
	GRPCClient
	// WSClient ws client
	WSClient
	// DB common db module
	DB
	// UserCustom User custom define named modules from 64-128
	UserCustom = 64
	// NamedMax named modules max ID
	ModuleMax = 255
)
