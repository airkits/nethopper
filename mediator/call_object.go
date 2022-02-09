package mediator

type ICaller interface {
	Execute(obj *CallObject) *RetObject
}

// CallObject call struct
type CallObject struct {
	Caller  ICaller
	CmdID   int32
	Option  int32
	Args    []interface{}
	ChanRet chan *RetObject
}

func (c *CallObject) Init(caller ICaller, cmdID int32, opt int32, args ...interface{}) *CallObject {
	c.Caller = caller
	c.CmdID = cmdID
	c.Option = opt
	c.Args = args
	c.ChanRet = make(chan *RetObject, 1)
	return c
}
func (c *CallObject) Reset() *CallObject {
	c.Args = nil
	c.Caller = nil
	c.CmdID = 0
	c.Option = 0
	if c.ChanRet != nil {
		close(c.ChanRet)
		c.ChanRet = nil
	}
	return c
}

type Callback func(ret *RetObject)

// RetObject call return object
type RetObject struct {
	Data interface{}
	Code int32
	Err  error
}

// NewCallObject create call object
func NewCallObject(caller ICaller, cmdID int32, opt int32, args ...interface{}) *CallObject {
	obj := &CallObject{}
	return obj.Init(caller, cmdID, opt, args...)
}

// NewRetObject create ret object
func NewRetObject(code int32, err error, data interface{}) *RetObject {
	obj := &RetObject{}
	return obj.Init(code, err, data)
}

func (c *RetObject) Init(code int32, err error, data interface{}) *RetObject {
	c.Code = code
	c.Data = data
	c.Err = err
	return c
}
func (c *RetObject) Reset() *RetObject {
	c.Code = 0
	c.Data = nil
	c.Err = nil
	return c
}
