package mediator

import "github.com/airkits/nethopper/config"

// MData is module info
type MData struct {
	ID         uint8
	CreateFunc func() (IModule, error)
	Conf       config.IConfig
	Priority   int32
	Module     IModule
}

func NewMData(mid uint8, createFunc func() (IModule, error), conf config.IConfig, dep []uint8) *MData {
	if dep == nil {
		dep = []uint8{}
	}
	priority := int32(mid)
	for _, e := range dep {
		priority += M().GetPriority(e)
	}
	data := &MData{
		ID:         mid,
		CreateFunc: createFunc,
		Conf:       conf,
		Priority:   priority,
	}
	return data
}

type MDataSlice []*MData

func (s MDataSlice) Len() int           { return len(s) }
func (s MDataSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s MDataSlice) Less(i, j int) bool { return s[i].Priority < s[j].Priority }
