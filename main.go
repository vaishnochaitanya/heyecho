package main

import (
	"vchan.in/heyecho/cmd/app"
	"vchan.in/heyecho/cmd/seed"
	"vchan.in/heyecho/db"
)

func main() {
	db.InitDB()
	seed.Seed()
	app.Server()
}
