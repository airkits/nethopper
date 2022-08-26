package base

import "errors"

// ErrQueueIsClosed queue is closed
var ErrQueueIsClosed = errors.New("Queue is Closed")

var ErrQueueEmpty = errors.New("Queue is empty")
var ErrQueueFull = errors.New("Queue is full")

// ErrQueueTimeout queue push or pop timeout
var ErrQueueTimeout = errors.New("Timeout on queue")

var ErrWriteChanTimeout = errors.New("write channel timeout")
var ErrReadChanTimeout = errors.New("read channel timeout")
var ErrCallTimeout = errors.New("Timeout on call")
