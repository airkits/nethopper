package common

const (
	// HeaderToken token key define
	HeaderToken = "token"
	// HeaderUID UID key define
	HeaderUID = "UID"
)

const (
	//PackageLengthSize package length size
	PackageLengthSize = 2
)

// ClientInfo grpc client info
type ClientInfo struct {
	ServerID int
	Name     string
	Address  string
}
