package rpc

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// eHTTPRequestType
const (
	// POST method
	POST = "POST"
	// GET method
	GET = "GET"
	//eHTTPRequestTypeText text format
	eHTTPRequestTypeText = 0x01
	//eHTTPRequestTypeJSON json format
	eHTTPRequestTypeJSON = 0x02
	//eHTTPRequestTypeByte byte buffer format
	eHTTPRequestTypeByte = 0x03
	//ContentTypeText request type text
	ContentTypeText = "application/x-www-form-urlencoded"
	//ContentTypeJSON request type json
	ContentTypeJSON = "application/json;charset=UTF-8"
	//ContentTypeByte request type buffer
	ContentTypeByte = "application/octet-stream" // "application/x-protobuf"  //
	//ResponseTypeText response type text
	ResponseTypeText = "text"
	//ResponseTypeJSON response type json
	ResponseTypeJSON = "json"
	//ResponseTypeXML response type xml
	ResponseTypeXML = "xml"
	//ResponseTypeByte response type buffer
	ResponseTypeByte = "arraybuffer"
)

const (
	//ConnTimeoutMS conn timeout
	ConnTimeoutMS = 1000
	//ServeTimeoutMS request timeout
	ServeTimeoutMS = 3000
)

// Request do http request with data
func Request(url string, method string, reqType int, header interface{}, data interface{}, responseType string, connTimeoutMs int, serveTimeoutMs int) (interface{}, error) {

	// client := &http.Client{
	// 	Transport: &http.Transport{
	// 		Dial: func(netw, addr string) (net.Conn, error) {
	// 			c, err := net.DialTimeout(netw, addr, time.Duration(connTimeoutMs)*time.Millisecond)
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 			c.SetDeadline(time.Now().Add(time.Duration(serveTimeoutMs) * time.Millisecond))
	// 			return c, nil
	// 		},
	// 	},
	// }
	// var body *bytes.Reader
	// var err error
	// if data != nil {
	// 	params := data.(map[string]interface{})
	// 	var bytesData []byte
	// 	switch reqType {
	// 	case eHTTPRequestTypeJSON:
	// 		bytesData, err = codec.JSONCodec.Marshal(params, nil)
	// 		break
	// 	}

	// 	body = bytes.NewReader(bytesData)
	// }

	// req, err := http.NewRequest(method, url, body)

	// if header != nil {
	// 	h := header.(map[string]interface{})
	// 	for key, val := range h {
	// 		req.Header.Set(key, val.(string))
	// 	}
	// }
	// contentType := ContentTypeText
	// switch reqType {
	// case eHTTPRequestTypeJSON:
	// 	contentType = ContentTypeJSON
	// 	break
	// case eHTTPRequestTypeText:
	// 	contentType = ContentTypeText
	// 	break
	// case eHTTPRequestTypeByte:
	// 	contentType = ContentTypeByte
	// }
	// req.Header.Set("Content-Type", contentType)
	// req.Response.Header.Set("")
	// resp, err := client.Do(req)
	// if err != nil {
	// 	err = fmt.Errorf("http failed, POST url:%s, reason:%s", url, err.Error())
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// defer resp.Body.Close()

	// if response.StatusCode != 200 {
	// 	err = fmt.Errorf("http status code error, POST url:%s, code:%d", url, response.StatusCode)
	// 	return
	// }

	// respBody, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	err = fmt.Errorf("cannot read http response, POST url:%s, reason:%s", url, err.Error())
	// 	return
	// }
	// str = string(respBody)

	return nil, nil
}

//HTTPGet get method
func HTTPGet(url string, connTimeoutMs int, serveTimeoutMs int) (str string, err error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Duration(connTimeoutMs)*time.Millisecond)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(time.Duration(serveTimeoutMs) * time.Millisecond))
				return c, nil
			},
		},
	}

	reqest, _ := http.NewRequest("GET", url, nil)
	response, err := client.Do(reqest)
	if err != nil {
		err = fmt.Errorf("http failed, GET url:%s, reason:%s", url, err.Error())
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = fmt.Errorf("http status code error, GET url:%s, code:%d", url, response.StatusCode)
		return
	}

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("cannot read http response, GET url:%s, reason:%s", url, err.Error())
		return
	}
	str = string(resBody)
	return
}

//HTTPPost post method
func HTTPPost(url string, data string, connTimeoutMs int, serveTimeoutMs int) (str string, err error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Duration(connTimeoutMs)*time.Millisecond)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(time.Duration(serveTimeoutMs) * time.Millisecond))
				return c, nil
			},
		},
	}

	body := strings.NewReader(data)
	reqest, err := http.NewRequest("POST", url, body)
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(reqest)
	if err != nil {
		err = fmt.Errorf("http failed, POST url:%s, reason:%s", url, err.Error())
		fmt.Println(err.Error())
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = fmt.Errorf("http status code error, POST url:%s, code:%d", url, response.StatusCode)
		return
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("cannot read http response, POST url:%s, reason:%s", url, err.Error())
		return
	}
	str = string(respBody)
	return

}
