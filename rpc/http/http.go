package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gonethopper/nethopper/codec"
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
	//eHTTPRequestTypePB protobuf3 byte buffer format
	eHTTPRequestTypePB = 0x03
	//eHTTPRequestTypeByte byte buffer format
	eHTTPRequestTypeByte = 0x04
	//ContentTypeText request type text
	ContentTypeText = "application/x-www-form-urlencoded"
	//ContentTypeJSON request type json
	ContentTypeJSON = "application/json;charset=UTF-8"
	//ContentTypeByte request type buffer
	ContentTypeByte = "application/octet-stream" // "application/x-protobuf"  //
	//ResponseTypeText response type text
	ResponseTypeText = 0x01
	//ResponseTypeJSON response type json
	ResponseTypeJSON = 0x02
	//ResponseTypeXML response type xml
	ResponseTypeXML = 0x03
	//ResponseTypePB response protobuf3 type buffer
	ResponseTypePB = 0x04
	//ResponseTypeByte response type buffer
	ResponseTypeByte = 0x05
)

const (
	//ConnTimeoutMS conn timeout
	ConnTimeoutMS = 1000
	//ServeTimeoutMS request timeout
	ServeTimeoutMS = 3000
)

// TextParamsValue get params formats
func TextParamsValue(param map[string]interface{}) []byte {
	var p = url.Values{}
	for key, value := range param {
		p.Add(key, value.(string))
	}
	return []byte(p.Encode())
}

// Request do http request with data
func Request(url string, method string, requestType int, header interface{}, data interface{}, responseType int, results interface{}, connTimeoutMs int, serveTimeoutMs int) error {

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
	var body *bytes.Reader
	var err error
	if data != nil {
		params := data.(map[string]interface{})
		var bytesData []byte
		switch requestType {
		case eHTTPRequestTypeJSON:
			{
				bytesData, err = codec.JSONCodec.Marshal(params, nil)
				break
			}
		case eHTTPRequestTypeText:
			{
				bytesData = TextParamsValue(params)
				break
			}
		case eHTTPRequestTypeByte:
			{
				bytesData, err = codec.BinaryCodec.Marshal(params, nil)
				break
			}
		case eHTTPRequestTypePB:
			{
				bytesData, err = codec.PBCodec.Marshal(params, nil)
				break
			}
		}
		if err != nil {
			return err
		}
		body = bytes.NewReader(bytesData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	if header != nil {
		h := header.(map[string]interface{})
		for key, val := range h {
			req.Header.Set(key, val.(string))
		}
	}
	contentType := ContentTypeText
	switch requestType {
	case eHTTPRequestTypeJSON:
		contentType = ContentTypeJSON
		break
	case eHTTPRequestTypeText:
		contentType = ContentTypeText
		break
	case eHTTPRequestTypePB:
	case eHTTPRequestTypeByte:
		contentType = ContentTypeByte
		break
	}
	req.Header.Set("Content-Type", contentType)
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("http failed, %s url:%s, reason:%s", method, url, err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("http status code error, %s url:%s, code:%d", method, url, resp.StatusCode)
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("cannot read http response, %s url:%s, reason:%s", method, url, err.Error())
		return err
	}

	switch responseType {
	case ResponseTypeJSON:
		{
			err = codec.JSONCodec.Unmarshal(respBody, results, nil)
			break
		}
	case ResponseTypeText:
		{
			results = string(respBody)
			break
		}
	case ResponseTypeByte:
		{
			err = codec.BinaryCodec.Unmarshal(respBody, results, nil)
			break
		}
	case ResponseTypePB:
		{
			err = codec.PBCodec.Unmarshal(respBody, results, nil)
			break
		}
	}
	return err
}
