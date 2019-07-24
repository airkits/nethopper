package connect

// Connect define network interface
type Connect interface {
	// Setup init Connect with config
	Setup(m map[string]interface{}) (Connect, error)
	// Listen on local port
	Listen()
	// Accept

}
