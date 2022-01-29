package log

// Log Levels Define
// 0       Fatal: system is unusable
// 1       Error: error conditions
// 2       Trace: trace conditions
// 3       Warning: warning conditions
// 4       Info: informational messages
// 5       Debug: debug-level messages
const (
	FATAL = iota
	ERROR
	TRACE
	WARNING
	INFO
	DEBUG
)
