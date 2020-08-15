package logic

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gonethopper/nethopper/examples/snowflake/cmd"
	"github.com/gonethopper/nethopper/log"
	"github.com/gonethopper/nethopper/network/common"
	"github.com/gonethopper/nethopper/server"
	"github.com/gonethopper/nethopper/utils"
)

func init() {
	conf := log.Config{
		Filename:     "logs/server_log.log",
		Level:        7,
		MaxLines:     1000,
		MaxSize:      50,
		HourEnabled:  true,
		DailyEnabled: true,
		QueueSize:    1000,
	}
	server.NewNamedModule(server.MIDLog, "log", log.LogModuleCreate, nil, &conf)
	logicConf := common.LogicConfig{
		QueueSize:          100,
		WorkerPoolCapacity: 100,
	}
	server.NewNamedModule(server.MIDLogic, "logic", ModuleCreate, nil, &logicConf)

}
func TestGetUID(t *testing.T) {
	v, result := server.Call(server.MIDLogic, cmd.CallIDGenUIDCmd, utils.RandomInt32(0, 1024), int32(1))
	t.Errorf("%v %v", v, result)
}
