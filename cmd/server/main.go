package main

import (
	"os"
	"tchat/server"
)

var (
	_ = os.Setenv("debug", "true")
)

func main() {
	server.Start()
}
