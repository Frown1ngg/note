package main

import (
	"notes_project/server"
)

func init() {
	server.InitServer()
}
func main() {
	server.StartServer()
}
