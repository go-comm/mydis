package mydis

import (
	"context"
)

func (db *DB) LPush(ctx context.Context, k []byte, v0 interface{}, v ...interface{}) error {
	return Update(ctx, db, func(tx Tx) error {
		return tx.LPush(ctx, k, v0, v...)
	})
}

func (db *DB) RPush(ctx context.Context, k []byte, v0 interface{}, v ...interface{}) error {
	return Update(ctx, db, func(tx Tx) error {
		return tx.RPush(ctx, k, v0, v...)
	})
}

func (db *DB) LPop(ctx context.Context, k []byte) (v interface{}, err error) {
	err = Update(ctx, db, func(tx Tx) error {
		v, err = tx.LPop(ctx, k)
		return err
	})
	return v, err
}

func (db *DB) RPop(ctx context.Context, k []byte) (v interface{}, err error) {
	err = Update(ctx, db, func(tx Tx) error {
		v, err = tx.RPop(ctx, k)
		return err
	})
	return v, err
}

func (db *DB) LLen(ctx context.Context, k []byte) (v interface{}, err error) {
	err = View(ctx, db, func(tx Tx) error {
		v, err = tx.LLen(ctx, k)
		return err
	})
	return
}

func (db *DB) LIndex(ctx context.Context, k []byte, i int) (v interface{}, err error) {
	err = View(ctx, db, func(tx Tx) error {
		v, err = tx.LIndex(ctx, k, i)
		return err
	})
	return
}

func (db *DB) LSet(ctx context.Context, k []byte, i int, v interface{}) error {
	return Update(ctx, db, func(tx Tx) error {
		return tx.LSet(ctx, k, i, v)
	})
}

func (db *DB) LRange(ctx context.Context, k []byte, start int, stop int) (v interface{}, err error) {
	err = View(ctx, db, func(tx Tx) error {
		v, err = tx.LRange(ctx, k, start, stop)
		return err
	})
	return v, err
}

func (db *DB) LTrim(ctx context.Context, k []byte, start int, stop int) error {
	return Update(ctx, db, func(tx Tx) error {
		return tx.LSet(ctx, k, start, stop)
	})
}
