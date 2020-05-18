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

//NodeInfo server node info
type NodeInfo struct {
	ID      int    `yarm:"id"`
	Name    string `yarm:"name"`
	Address string `yarm:"address"`
}
