package db

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gonethopper/nethopper/database"
	"github.com/gonethopper/nethopper/examples/usercenter/cmd"
	"github.com/gonethopper/nethopper/log"
	"github.com/gonethopper/nethopper/server"
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
	info := database.NodeInfo{
		ID:     1,
		Driver: "mysql",
		DSN:    "root:123456@tcp(127.0.0.1:3306)/usercenter?charset=utf8&parseTime=True&loc=Asia%2FShanghai",
	}
	dbConf := database.Config{
		Nodes:           []database.NodeInfo{info},
		ConnectInterval: 10,
		QueueSize:       100,
	}
	server.NewNamedModule(server.MIDDB, "mysql", ModuleCreate, nil, &dbConf)

}
func TestGetUserTableByUID(t *testing.T) {
	tn := getUserTableByUID(1112)
	if tn != "usercenter.user_0" {
		t.Error(tn)
	}

}
func TestGetOID2UIDTable(t *testing.T) {
	tn := getOID2UIDTable("openID")
	if tn != "usercenter.oid2uid_4" {
		t.Error(tn)
	}
}
func TestGetUIDByOpenID(t *testing.T) {
	openID := "openID"
	result, err := server.Call(server.MIDDB, cmd.MCDBGetUIDByOpenID, 0, openID)
	if err != nil {
		t.Error(err.Error())
		return
	}
	server.Info("get uid %d", result.(uint64))

}

func TestInsertOID2UID(t *testing.T) {
	openID := "openID"
	uid := uint64(1123456)
	result, err := server.Call(server.MIDDB, cmd.MCDBInsertOID2UID, 0, openID, uid)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("get uid %v", result.(bool))
}
