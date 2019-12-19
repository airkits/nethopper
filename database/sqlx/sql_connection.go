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
// OUT OF OR IN SQLConnection WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// * @Author: ankye
// * @Date: 2019-07-16 21:53:29
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-07-16 21:53:29

package sqlx

import (
	"database/sql"

	"github.com/gonethopper/nethopper/server"
	"github.com/jmoiron/sqlx"
)

// NewSQLConnection create redis cache instance
func NewSQLConnection(m map[string]interface{}) (*SQLConnection, error) {
	conn := &SQLConnection{}
	return conn.Setup(m)

}

// SQLConnection connect to db by dsn
type SQLConnection struct {
	db     *sqlx.DB
	Driver string
	DSN    string
}

// Setup init cache with config
func (s *SQLConnection) Setup(m map[string]interface{}) (*SQLConnection, error) {
	if err := s.ReadConfig(m); err != nil {
		return nil, err
	}
	return s, nil
}

// ReadConfig config map
// driver default mysql
// dsn default "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Asia%2FShanghai"
func (s *SQLConnection) ReadConfig(m map[string]interface{}) error {

	if err := server.ParseConfigValue(m, "driver", "mysql", &s.Driver); err != nil {
		return err
	}

	if err := server.ParseConfigValue(m, "dsn", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Asia%2FShanghai", &s.DSN); err != nil {
		return err
	}

	return nil
}

//Open connect and ping
func (s *SQLConnection) Open() error {
	var err error
	if s.db, err = sqlx.Connect(s.Driver, s.DSN); err != nil {
		panic(err.Error())
	}
	return s.Ping()
}

//Ping test SQLConnection
func (s *SQLConnection) Ping() error {
	// force a SQLConnection and test ping
	err := s.db.Ping()
	if err != nil {
		server.Error("couldn't connect to database: %s %s", s.Driver, s.DSN)
		panic(err.Error())
	}
	return err
}

//Close close SQLConnection
func (s *SQLConnection) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

//IsErrNoRows 判断是否有数据
func (s *SQLConnection) IsErrNoRows(err error) bool {
	return sql.ErrNoRows == err
}

//Select select operate
func (s *SQLConnection) Select(dest interface{}, query string, args ...interface{}) error {
	return s.db.Select(dest, query, args...)
}

//Exec process sql and get result
func (s *SQLConnection) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

//QueryRow by sql
func (s *SQLConnection) QueryRow(query string, args ...interface{}) *sqlx.Row {
	return s.db.QueryRowx(query, args...)
}

//Query sql and return rows
func (s *SQLConnection) Query(query string, args ...interface{}) (*sqlx.Rows, error) {
	return s.db.Queryx(query, args...)
}
