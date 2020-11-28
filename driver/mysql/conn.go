package mysql

import (
	"context"
	"database/sql"

	"github.com/go-comm/mydis"
)

const (
	createSQL = `CREATE TABLE %s (
	k varchar(127) NOT NULL,
	v varchar(255) DEFAULT NULL,
	ex bigint(20) DEFAULT NULL,
	ctime bigint(20) DEFAULT NULL,
	PRIMARY KEY (k)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
`
)

type mysqlConn struct {
	DB        *sql.DB
	tableName string
	getSQL    string
	setSQL    string
	delSQL    string
}

func (c *mysqlConn) handleExpiresInLoop() {

}

func (c *mysqlConn) Begin(ctx context.Context, readOnly bool) (mydis.Tx, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &mysqlTx{Tx: tx}, nil
}

func (c *mysqlConn) Get(ctx context.Context, key []byte) (v interface{}, err error) {
	row := c.DB.QueryRowContext(ctx, c.getSQL, mydis.BytesToString(key))
	var d []byte
	err = row.Scan(&d)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, mydis.ErrNoKey
		}
		return nil, err
	}
	return d, err
}

func (c *mysqlConn) Set(ctx context.Context, key []byte, v interface{}) (err error) {
	_, err = c.DB.ExecContext(ctx, c.setSQL, mydis.BytesToString(key), v, -1)
	return
}
