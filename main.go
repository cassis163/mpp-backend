package main

import (
	"mpp/database"
	"mpp/routing"
)

func main() {
	db := database.Init()
	routing.Init(db)
	defer db.Close()
}
