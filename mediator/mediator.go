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
	m.ref = base.NewRef()
	return m
}

type Mediator struct {
	modules [ModuleMax]IModule //module id => MData
	datas   MDataSlice         // array mdata cache
	wg      sync.WaitGroup
	ref     base.IRef
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
	base.GOFunctionWithWG(m.wg, m.ref, ModuleRun, module)
	return module, nil
}

func (m *Mediator) GetModuleByID(mid uint8) IModule {
	return m.modules[mid]
}

// func (m *Mediator) createModuleByID(MID int32, name string, parent IModule, conf config.IConfig) (IModule, error) {
// 	m, err := CreateModule(name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	m.MakeContext(nil, int32(conf.GetQueueSize()))
// 	m.SetName(ModuleName(m))
// 	cmdRegister(m)
// 	m.Setup(conf)
// 	m.SetID(MID)
// 	App.Modules.Store(MID, m)
// 	if MID == MIDLog {
// 		GLoggerModule = m
// 	}
// 	GOWithContext(ModuleRun, m)
// 	return m, nil
// }

// // CreateModule create module by name
// func (m *Mediator) CreateModule(name string) (IModule, error) {
// 	if f, ok := relModules[name]; ok {
// 		return f()
// 	}
// 	return nil, fmt.Errorf("You need register Module %s first", name)
// }

// // GetModuleByID get module instance by id
// func (m *Mediator) GetModuleByID(MID int32) (IModule, error) {
// 	se, ok := App.Modules.Load(MID)
// 	if ok {
// 		return se.(Module), nil
// 	}
// 	return nil, fmt.Errorf("cant get module ID %d", MID)
// }

// // NewModule create anonymous module
// func NewModule(name string, parent IModule, conf config.IConfig) (IModule, error) {
// 	//Inc AnonymousMID count = count +1
// 	MID := atomic.AddInt32(&AnonymousMID, 1)
// 	return createModuleByID(MID, name, parent, conf)
// }
