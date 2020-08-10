package sdk

import (
	"errors"
	"fmt"

	"github.com/gonethopper/nethopper/examples/usercenter/model"
	"github.com/gonethopper/nethopper/network/http"
	"github.com/gonethopper/nethopper/server"
)

//GenUIDAPI api server
const GenUIDAPI = "%s/v1/genuid"

//GenUID 获取uid
func GenUID(host string, channel int32) (uint64, error) {
	url := fmt.Sprintf(GenUIDAPI, host)

	req := map[string]interface{}{
		"Channel": channel,
	}
	resp := model.Response{
		Data: &model.GenUIDResp{},
	}
	if err := http.Request(url, http.POST, http.RequestTypeJSON, nil, req, http.ResponseTypeJSON, &resp, http.ConnTimeoutMS, http.ServeTimeoutMS); err != nil {
		return 0, err
	}
	server.Info("%v", resp)
	if resp.Code == 0 {
		return resp.Data.(*model.GenUIDResp).UID, nil
	}
	return 0, errors.New(resp.Msg)

}
