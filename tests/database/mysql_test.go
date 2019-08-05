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
// * @Date: 2019-06-21 13:42:42
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-21 13:42:42

package database_test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gonethopper/nethopper/database/sqlx"
	. "github.com/gonethopper/nethopper/server"
)

func TestSQLConnection(t *testing.T) {
	m := map[string]interface{}{
		"queueSize": 1000,
		"dsn":       "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Asia%2FShanghai",
		"driver":    "mysql",
	}
	if conn, err := sqlx.NewSQLConnection(m); err == nil {
		if err := conn.Open(); err != nil {
			t.Error(err)
			Error("error")
		}
	} else {
		t.Error(err)
	}

}
