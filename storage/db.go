package storage

import (
	"github.com/jackc/pgx"
	"github.com/aspcartman/eros/env"
	"github.com/jackc/pgx/log/logrusadapter"
	"github.com/aspcartman/exceptions"
	"fmt"
	"github.com/sirupsen/logrus"
	"unsafe"
)

type DB struct {
	pool *pgx.ConnPool
}

func NewDB(host string, port uint16, user, pass, database string) DB {
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     host,
			Port:     port,
			User:     user,
			Password: pass,
			Database: database,
			Logger:   logrusadapter.NewLogger(env.Log),
			LogLevel: pgx.LogLevelWarn,
		},
	})
	if err != nil {
		e.Throw("failed creating a connection pool", err, e.Map{
			"Host":     host,
			"Port":     port,
			"User":     user,
			"Password": pass,
			"Database": database,
		})
	}

	db := DB{pool}
	return db
}

func (d *DB) Upgrade() {
	tx := d.Begin()
	defer tx.Rollback()

	// todo lock

	current := d.currentRevision(tx)
	if current == len(revisions)-1 {
		env.Log.WithField("revision", current).Info("DB is up to the current revision, no migrations needed")
		return
	}

	env.Log.WithFields(logrus.Fields{
		"current": current,
		"target":  len(revisions) - 1,
	}).Info("DB requires upgrade. Prefoming migrations")

	for i := current + 1; i < len(revisions); i++ {
		env.Log.WithField("revision", i).Info("Upgrading")

		revisions[i](tx)

		tx.Exec(`UPDATE app SET value = $1 WHERE key='revision'`, fmt.Sprint(i))
	}

	tx.Commit()

	env.Log.WithField("revision", len(revisions)-1).Info("Succesfully upgraded")
}

func (d *DB) currentRevision(tx *Tx) int {
	if !d.exists(tx, "app") {
		return -1
	}

	rev := 0
	err := tx.QueryRow(`SELECT value FROM app WHERE key='revision'`).Scan(&rev)
	if err != nil {
		e.Throw("failed getting current db revision", err)
	}
	return rev
}

func (d *DB) exists(tx *Tx, table string) bool {
	exists := false
	err := tx.QueryRow(`SELECT EXISTS (
	   SELECT 1
	   FROM   information_schema.tables
	   WHERE  table_schema = 'public'
	   AND    table_name = $1
    )`, table).Scan(&exists)
	if err != nil {
		e.Throw("failed checking for table existance", err)
	}
	return exists
}

func (d *DB) Begin() *Tx {
	return (*Tx)(unsafe.Pointer(d.tx()))
}

func (d *DB) tx() *pgx.Tx {
	tx, err := d.pool.Begin()
	if err != nil {
		e.Throw("failed starting bd transaction", err)
	}
	return tx
}
