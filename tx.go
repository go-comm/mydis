package mydis

import (
	"context"
)

type Tx interface {
	Commit() error
	Rollback() error
	Get(ctx context.Context, k []byte) (interface{}, error)
	Set(ctx context.Context, k []byte, v interface{}) error
	Del(ctx context.Context, k []byte) (err error)

	// list
	LPush(ctx context.Context, k []byte, v0 interface{}, v ...interface{}) error
	RPush(ctx context.Context, k []byte, v0 interface{}, v ...interface{}) error
	LPop(ctx context.Context, k []byte) (interface{}, error)
	RPop(ctx context.Context, k []byte) (interface{}, error)
	LLen(ctx context.Context, k []byte) (interface{}, error)
	LIndex(ctx context.Context, k []byte, i int) (interface{}, error)
	LSet(ctx context.Context, k []byte, i int, v interface{}) error
	LRange(ctx context.Context, k []byte, start int, stop int) (interface{}, error)
	LTrim(ctx context.Context, k []byte, start int, stop int) error
}

func View(ctx context.Context, db *DB, fn func(tx Tx) error) error {
	tx, err := db.Begin(ctx, true)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		return tx.Rollback()
	}
	return tx.Commit()
}

func Update(ctx context.Context, db *DB, fn func(tx Tx) error) error {
	tx, err := db.Begin(ctx, false)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		return tx.Rollback()
	}
	return tx.Commit()
}
