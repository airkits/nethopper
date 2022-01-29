package mq

import "sync"

// GBytesPool pre-create []byte pool
var GBytesPool *BytesPool

// GMessagePool pre-create message pool
var GMessagePool *MessagePool

var once sync.Once

//GMQ session instance
func GMQ() *MessagePool {
	once.Do(func() {
		GMessagePool = NewMessagePool()
		GBytesPool = NewBytesPool()
	})
	return GMessagePool
}
