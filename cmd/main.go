package main

import (
	"flag"
	"strings"

	"github.com/justxuewei/ip-monitor/pkg"
)

func main() {
	sendKey := flag.String("key", "", "SendKey for ServerChan")
	serverName := flag.String("name", "IPMONITOR", "Server name")
	devicesStr := flag.String("devices", "", "Devices you want to monitor, e.g. \"lo,enp5s0\"")
	heartbeat := flag.Bool("heartbeat", false, "Sending heartbeat everyday to tell you your machine is alive")

	flag.Parse()

	if *sendKey == "" {
		panic("SendKey is required")
	}

	var devices []string
	if *devicesStr != "" {
		devices = strings.Split(*devicesStr, ",")
	}

	pusher := pkg.NewServerChan(*sendKey)
	monitor := pkg.NewMonitor(*serverName, pusher, devices, *heartbeat)

	monitor.Check()
}
