// MIT License

// Copyright (c) 2019 gonethopper

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// * @Author: ankye
// * @Date: 2019-12-11 10:13:21
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-12-11 10:13:21

package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/airkits/nethopper/codec"
)

//ConnSet http conn set
type ConnSet map[net.Conn]struct{}

// eHTTPRequestType
const (
	// POST method
	POST = "POST"
	// GET method
	GET = "GET"

	//DELETE method
	DELETE = "DELETE"

	//RequestTypeText text format
	RequestTypeText = 0x01
	//RequestTypeJSON json format
	RequestTypeJSON = 0x02
	//RequestTypePB protobuf3 byte buffer format
	RequestTypePB = 0x03
	//RequestTypeByte byte buffer format
	RequestTypeByte = 0x04
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

	//ConnTimeoutMS conn timeout
	ConnTimeoutMS = 5000
	//ServeTimeoutMS request timeout
	ServeTimeoutMS = 10000
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
	var bindIP string
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				lAddr, err := net.ResolveTCPAddr(netw, bindIP+":0")
				if err != nil {
					return nil, err
				}
				// //被请求的地址
				// rAddr, err := net.ResolveTCPAddr(netw, addr)
				// if err != nil {
				// 	return nil, err
				// }
				// c, err := net.DialTCP(netw, lAddr, rAddr)
				// if err != nil {
				// 	return nil, err
				// }

				d := net.Dialer{Timeout: time.Duration(connTimeoutMs) * time.Millisecond, LocalAddr: lAddr}
				c, err := d.Dial("tcp", addr)
				if err != nil {
					return nil, err
				}

				// c, err := net.DialTimeout(netw, addr, time.Duration(connTimeoutMs)*time.Microsecond)
				// if err != nil {
				// 	return nil, err
				// }

				c.SetDeadline(time.Now().Add(time.Duration(serveTimeoutMs) * time.Millisecond))
				return c, nil
			},
		},
	}
	var body *bytes.Reader
	var err error
	var req *http.Request
	if data != nil {
		params := data.(map[string]interface{})
		var bytesData []byte
		switch requestType {
		case RequestTypeJSON:
			{
				bytesData, err = codec.JSONCodec.Marshal(params)
				break
			}
		case RequestTypeText:
			{
				bytesData = TextParamsValue(params)
				break
			}
		case RequestTypeByte:
			{
				bytesData, err = codec.RawCodec.Marshal(params)
				break
			}
		case RequestTypePB:
			{
				bytesData, err = codec.PBCodec.Marshal(params)
				break
			}
		}
		if err != nil {
			return err
		}
		body = bytes.NewReader(bytesData)
		req, err = http.NewRequest(method, url, body)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return err
	}
	// parnet := ctx
	// if parnet == nil {
	// 	parnet = context.TODO()
	// }
	// currentCtx, cancel := context.WithCancel(parnet)
	// time.AfterFunc(time.Duration(serveTimeoutMs)*time.Millisecond, func() {
	// 	cancel()
	// })
	// req = req.WithContext(currentCtx)

	if header != nil {
		h := header.(map[string]interface{})
		for key, val := range h {
			req.Header.Set(key, val.(string))
		}
	}
	contentType := ContentTypeText
	switch requestType {
	case RequestTypeJSON:
		contentType = ContentTypeJSON
		break
	case RequestTypeText:
		contentType = ContentTypeText
		break
	case RequestTypePB:
	case RequestTypeByte:
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
			err = codec.JSONCodec.Unmarshal(respBody, results)
			break
		}
	case ResponseTypeText:
		{
			*(results.(*string)) = string(respBody)
			break
		}
	case ResponseTypeByte:
		{
			err = codec.RawCodec.Unmarshal(respBody, results)
			break
		}
	case ResponseTypePB:
		{
			err = codec.PBCodec.Unmarshal(respBody, results)
			break
		}
	}
	return err
}
