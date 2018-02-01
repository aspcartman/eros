package storage

import (
	"github.com/jackc/pgx"
	"github.com/aspcartman/exceptions"
	"unsafe"
)

type Tx struct {
	pgx.Tx
}

func (tx *Tx) Exec(sql string, arguments ...interface{}) uint {
	cmd, err := tx.Tx.Exec(sql, arguments...)
	if err != nil {
		e.Throw("failure executing statement", err, e.Map{
			"sql":  sql,
			"args": arguments,
		})
	}
	return uint(cmd.RowsAffected())
}

func (tx *Tx) Commit() {
	err := tx.Tx.Commit()
	if err != nil {
		e.Throw("commit failed", err)
	}
}

func (tx *Tx) Query(sql string, args ...interface{}) *Rows {
	rows, err := tx.Tx.Query(sql, args...)
	if err != nil {
		e.Throw("failed making query", err, e.Map{
			"query": sql,
			"args":  args,
		})
	}
	return (*Rows)(unsafe.Pointer(rows))
}
