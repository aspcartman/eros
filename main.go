package main

import (
	"github.com/aspcartman/eros/storage"
)

func main() {
	db := storage.NewDB("localhost", 7771, "postgres", "", "eros")
	db.Upgrade()

}