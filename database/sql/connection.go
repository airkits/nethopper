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
// * @Date: 2019-07-16 21:53:29
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-07-16 21:53:29

package sql

import (
	"database/sql"

	"github.com/gonethopper/nethopper/server"
	"github.com/jmoiron/sqlx"
)

// NewConnection create redis cache instance
func NewConnection(m map[string]interface{}) (*Connection, error) {
	conn := &Connection{}
	return conn.Setup(m)

}

// Connection connect to db by dsn
type Connection struct {
	db     *sqlx.DB
	Driver string
	DSN    string
}

// Setup init cache with config
func (s *Connection) Setup(m map[string]interface{}) (*Connection, error) {
	if err := s.readConfig(m); err != nil {
		return nil, err
	}
	return s, nil
}

// config map
// driver default mysql
// dsn default "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Asia%2FShanghai"
func (s *Connection) readConfig(m map[string]interface{}) error {

	driver, err := server.ParseValue(m, "driver", "mysql")
	if err != nil {
		return err
	}
	s.Driver = driver.(string)
	dsn, err := server.ParseValue(m, "dsn", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Asia%2FShanghai")
	if err != nil {
		return err
	}
	s.DSN = dsn.(string)

	return nil
}

//Open connect and ping
func (s *Connection) Open() error {
	var err error
	if s.db, err = sqlx.Connect(s.Driver, s.DSN); err != nil {
		panic(err.Error())
	}
	return s.Ping()
}

//Ping test connection
func (s *Connection) Ping() error {
	// force a connection and test ping
	err := s.db.Ping()
	if err != nil {
		server.Error("couldn't connect to database: %s %s", s.Driver, s.DSN)
		panic(err.Error())
	}
	return err
}

//Close close connection
func (s *Connection) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

//IsErrNoRows 判断是否有数据
func (s *Connection) IsErrNoRows(err error) bool {
	return sql.ErrNoRows == err
}

//Select select operate
func (s *Connection) Select(dest interface{}, query string, args ...interface{}) error {
	return s.db.Select(dest, query, args...)
}

//Exec process sql and get result
func (s *Connection) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

//QueryRow by sql
func (s *Connection) QueryRow(query string, args ...interface{}) *sqlx.Row {
	return s.db.QueryRowx(query, args...)
}

//Query sql and return rows
func (s *Connection) Query(query string, args ...interface{}) (*sqlx.Rows, error) {
	return s.db.Queryx(query, args...)
}
