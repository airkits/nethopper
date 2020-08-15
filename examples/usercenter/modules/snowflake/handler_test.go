package snowflake

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gonethopper/nethopper/examples/usercenter/cmd"
	"github.com/gonethopper/nethopper/examples/usercenter/global"
	"github.com/gonethopper/nethopper/examples/usercenter/model"
	"github.com/gonethopper/nethopper/log"
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
	sfConf := model.SFConfig{
		Hosts:     []string{"http://127.0.0.1:11080"},
		QueueSize: 100,
	}
	server.NewNamedModule(global.MIDSnowflakeClient, "sfClient", ModuleCreate, nil, &sfConf)

}
func TestGetUID(t *testing.T) {
	v, result := server.Call(global.MIDSnowflakeClient, cmd.SFGetUID, utils.RandomInt32(0, 1024), int32(1))
	t.Errorf("%v %v", v, result)
}
