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
	LPush(ctx context.Context, k []byte, v interface{}) error
	RPush(ctx context.Context, k []byte, v interface{}) error
	LPop(ctx context.Context, k []byte) (interface{}, error)
	RPop(ctx context.Context, k []byte) (interface{}, error)
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
