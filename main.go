package main

import (
	"github.com/aspcartman/eros/storage"
)

func main() {
	db := storage.NewDB("postgres", 5432, "postgres", "", "eros")
	db.Upgrade()

}
