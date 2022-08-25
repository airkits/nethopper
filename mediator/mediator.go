package mediator

import (
	"sort"
	"sync"

	"github.com/airkits/nethopper/base"
)

func NewMediator() *Mediator {
	m := new(Mediator)
	m.modules = [ModuleMax]IModule{}
	m.datas = make(MDataSlice, 0)
	m.AppCtx = base.NewAppContext()
	return m
}

type Mediator struct {
	modules [ModuleMax]IModule //module id => MData
	datas   MDataSlice         // array mdata cache
	AppCtx  *base.AppContext
	sync.Mutex
}

func (m *Mediator) GetPriority(mid uint8) int32 {
	m.Lock()
	defer m.Unlock()
	for _, e := range m.datas {
		if e.ID == mid {
			return e.Priority
		}
	}
	return 0
}

func (m *Mediator) HasModule(mid uint8) bool {
	m.Lock()
	defer m.Unlock()
	for _, e := range m.datas {
		if e.ID == mid {
			return true
		}
	}
	return false
}

func (m *Mediator) CreateModule(data *MData) (IModule, error) {
	m.Lock()
	defer m.Unlock()
	module, err := data.CreateFunc()
	if err != nil {
		return nil, err
	}

	module.MakeContext(int32(data.Conf.GetQueueSize()))
	module.Setup(data.Conf)
	module.SetID(data.ID)
	module.SetName(ModuleName(module))
	cmdRegister(module)
	data.Module = module
	m.modules[data.ID] = module
	m.datas = append(m.datas, data)

	sort.Sort(m.datas)

	base.GOFunctionWithWG(m.AppCtx.WG, m.AppCtx.Ref, ModuleRun, module)
	return module, nil
}
func (m *Mediator) Wait() {
	m.AppCtx.WG.Wait()
}
func (m *Mediator) Exit() {
	for i := ModuleMax - 1; i >= 0; i-- {
		m := m.GetModuleByID(uint8(i))
		if m != nil {
			m.Close()
		}
	}
}
func (m *Mediator) GetModuleByID(mid uint8) IModule {
	return m.modules[mid]
}
