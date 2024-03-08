package dbtx

import (
	"context"
	"errors"

	"github.com/fjrid/parking/pkg/dbtxn"
)

type dbTx struct {
	tx *dbtxn.Context
}

// Begin used to start transaction
func Begin(ctx *context.Context) *dbTx {
	tx := dbtxn.Begin(ctx)
	return &dbTx{
		tx: tx,
	}
}

// Commit used to commit transaction
func (d *dbTx) Commit() (err error) {
	if d.tx != nil {
		err = d.tx.Commit()
		d.tx = nil
	}
	return
}

// Rollback used to rollback transaction
func (d *dbTx) Rollback(e error) (err error) {
	if d.tx != nil {
		if e == nil {
			e = errors.New("rolled back")
		}

		d.tx.AppendError(e)
		err = d.tx.Commit()
		d.tx = nil
	}
	return
}
