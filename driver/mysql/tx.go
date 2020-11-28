package mysql

import (
	"context"
	"database/sql"

	"github.com/go-comm/mydis"
)

type mysqlTx struct {
	Tx *sql.Tx
	c  *mysqlConn
}

func (tx *mysqlTx) Commit() error {
	return tx.Commit()
}

func (tx *mysqlTx) Rollback() error {
	return tx.Rollback()
}

func (tx *mysqlTx) Scan(ctx context.Context, dest interface{}, k []byte) (err error) {
	row := tx.Tx.QueryRowContext(ctx, tx.c.getSQL, mydis.BytesToString(k))
	err = row.Scan(dest)
	if err != nil {
		if err == sql.ErrNoRows {
			return mydis.ErrNoKey
		}
		return err
	}
	return err
}

func (tx *mysqlTx) Get(ctx context.Context, k []byte) (v interface{}, err error) {
	row := tx.Tx.QueryRowContext(ctx, tx.c.getSQL, mydis.BytesToString(k))
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

func (tx *mysqlTx) Set(ctx context.Context, k []byte, v interface{}) (err error) {
	_, err = tx.Tx.ExecContext(ctx, tx.c.setSQL, mydis.BytesToString(k), v, -1)
	return
}

func (tx *mysqlTx) Del(ctx context.Context, k []byte) (err error) {
	_, err = tx.Tx.ExecContext(ctx, tx.c.delSQL, mydis.BytesToString(k))
	return
}
