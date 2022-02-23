package base

import "sync"

type AppContext struct {
	WG  *sync.WaitGroup
	Ref IRef
}

func NewAppContext() *AppContext {
	ctx := &AppContext{}
	ctx.Ref = NewRef()
	ctx.WG = &sync.WaitGroup{}
	return ctx
}
