package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/justxuewei/ip-monitor/pkg"
)

var version = "dev"

func main() {
	webhookURL := flag.String("webhook-url", "", "Webhook URL template containing {message}")
	serverName := flag.String("name", "IPMONITOR", "Server name")
	devicesStr := flag.String("devices", "", "Devices you want to monitor, e.g. \"lo,enp5s0\"")
	heartbeat := flag.Bool("heartbeat", false, "Sending heartbeat everyday to tell you your machine is alive")
	showVersion := flag.Bool("version", false, "Print version and exit")

	flag.CommandLine.SetOutput(os.Stdout)
	flag.Usage = printUsage

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help":
			printUsage()
			return
		case "version":
			printVersion()
			return
		}
	}

	flag.Parse()

	if *showVersion {
		printVersion()
		return
	}

	if *webhookURL == "" {
		panic("webhook-url is required")
	}

	var devices []string
	if *devicesStr != "" {
		devices = strings.Split(*devicesStr, ",")
	}

	pusher := pkg.NewServerChan(*webhookURL)
	monitor := pkg.NewMonitor(*serverName, pusher, devices, *heartbeat)

	monitor.Check()
}

func printVersion() {
	fmt.Println(version)
}

func printUsage() {
	fmt.Fprintln(flag.CommandLine.Output(), "Usage:")
	fmt.Fprintln(flag.CommandLine.Output(), "  ipmonitor [options]")
	fmt.Fprintln(flag.CommandLine.Output(), "  ipmonitor help")
	fmt.Fprintln(flag.CommandLine.Output(), "  ipmonitor version")
	fmt.Fprintln(flag.CommandLine.Output())
	fmt.Fprintln(flag.CommandLine.Output(), "Options:")
	flag.PrintDefaults()
}
