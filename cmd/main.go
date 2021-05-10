package main

import (
	"github.com/Utsavk/capture-events/cmd/server"
	"github.com/Utsavk/capture-events/config"
)

func main() {
	if !config.Parse() {
		return
	}
	server.StartServer(config.Props.Server)
}
