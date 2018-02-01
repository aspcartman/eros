package storage

import (
	"github.com/jackc/pgx"
	"github.com/aspcartman/pcache/e"
)

type Rows struct {
	pgx.Rows
}

func (r *Rows) Scan(dest ...interface{}) {
	err := r.Rows.Scan(dest...)
	if err != nil {
		e.Throw("scanning failed", err)
	}
}
