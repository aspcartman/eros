package main

import (
	"github.com/aspcartman/eros/storage"
)

func main() {
	db := storage.NewDB("postgres", 5432, "postgres", "", "eros")
	db.Upgrade()

	tx1, tx2 := db.Begin(), db.Begin()

	for p := range tx1.AllPhotos() {
		p.Refresh()
		tx2.SavePhoto(&p)
	}

	tx2.Commit()
}
