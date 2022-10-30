package main

import (
	"flag"

	"github.com/justxuewei/ip-monitor/pkg"
)

func main() {
	sendKey := flag.String("key", "", "SendKey for ServerChan")
	serverName := flag.String("name", "IPMONITOR", "Server Name")

	flag.Parse()

	if *sendKey == "" {
		panic("SendKey is required")
	}

	pusher := pkg.NewServerChan(*sendKey)
	monitor := pkg.NewMonitor(*serverName, pusher)

	monitor.Check()
}
