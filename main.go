package main

import (
	"vchan.in/heyecho/cmd/app"
	"vchan.in/heyecho/db"
)

func main() {
	db.InitDB()
	app.Server()
}
