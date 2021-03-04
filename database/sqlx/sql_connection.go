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

	"github.com/airkits/nethopper/database"
	"github.com/airkits/nethopper/server"
	"github.com/jmoiron/sqlx"
)

// NewSQLConnection create redis cache instance
func NewSQLConnection(nodes []database.NodeInfo) (*SQLConnection, error) {
	conn := &SQLConnection{}
	return conn.Setup(nodes)

}

// SQLConnection connect to db by dsn
type SQLConnection struct {
	pools []*sqlx.DB
	Nodes []database.NodeInfo
}

// Setup init cache with config
func (s *SQLConnection) Setup(nodes []database.NodeInfo) (*SQLConnection, error) {
	s.Nodes = nodes
	s.pools = make([]*sqlx.DB, len(nodes))
	return s, nil
}

//Open connect and ping
func (s *SQLConnection) Open() error {
	for index, info := range s.Nodes {
		db, err := sqlx.Connect(info.Driver, info.DSN)
		if err != nil {
			panic(err.Error())
		}
		s.pools[index] = db
	}

	return s.Ping()
}

//Ping test SQLConnection
func (s *SQLConnection) Ping() error {
	// force a SQLConnection and test ping
	for index, db := range s.pools {
		err := db.Ping()
		if err != nil {
			server.Error("couldn't connect to database: %s %s", s.Nodes[index].Driver, s.Nodes[index].DSN)
			panic(err.Error())
		}
		return err
	}
	return nil
}

//Close close SQLConnection
func (s *SQLConnection) Close() {
	for index, db := range s.pools {
		if db != nil {
			db.Close()
			server.Info("close db connection: %s %s", s.Nodes[index].Driver, s.Nodes[index].DSN)
		}
	}

}
func (s *SQLConnection) db() *sqlx.DB {
	return s.pools[0]
}

//IsErrNoRows 判断是否有数据
func (s *SQLConnection) IsErrNoRows(err error) bool {
	return sql.ErrNoRows == err
}

//Select select operate
func (s *SQLConnection) Select(dest interface{}, query string, args ...interface{}) error {
	return s.db().Select(dest, query, args...)
}

//Get get struct
func (s *SQLConnection) Get(dest interface{}, query string, args ...interface{}) error {
	return s.db().Get(dest, query, args...)
}

//Exec process sql and get result
func (s *SQLConnection) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := s.db().Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

//QueryRow by sql
func (s *SQLConnection) QueryRow(query string, args ...interface{}) *sqlx.Row {
	return s.db().QueryRowx(query, args...)
}

//Query sql and return rows
func (s *SQLConnection) Query(query string, args ...interface{}) (*sqlx.Rows, error) {
	return s.db().Queryx(query, args...)
}

//CreateTx create sqlx tx
func (s *SQLConnection) CreateTx() (*sqlx.Tx, error) {
	return s.db().Beginx()
}

//ExecTX exec tx
func (s *SQLConnection) ExecTX(querys []string) error {
	tx := s.db().MustBegin()
	for _, sql := range querys {
		tx.MustExec(sql)
	}
	return tx.Commit()
}
