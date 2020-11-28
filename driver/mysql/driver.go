package mysql

import (
	"database/sql"

	"github.com/go-comm/mydis"
)

type MysqlDriver struct {
}

func (d *MysqlDriver) Open(dsn string) (mydis.Conn, error) {
	return NewConn(sql.Open("mysql", dsn))
}

func NewConn(db *sql.DB, err error) (mydis.Conn, error) {
	if err != nil {
		return nil, err
	}
	c := &mysqlConn{DB: db}
	go c.handleExpiresInLoop()
	return c, nil
}
