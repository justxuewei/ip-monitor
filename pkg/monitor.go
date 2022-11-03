package pkg

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/vishvananda/netlink"
)

const (
	netInfoFile   = "/tmp/ip-monitor.txt"
	heartbeatFile = "/tmp/ip-monitor.updated-at.txt"
)

type Monitor struct {
	name      string
	devices   []string
	heartbeat bool

	pusher *ServerChan
}

func NewMonitor(name string, pusher *ServerChan, devices []string, heartbeat bool) *Monitor {
	return &Monitor{
		name:      name,
		devices:   devices,
		heartbeat: heartbeat,
		pusher:    pusher,
	}
}

func (m *Monitor) doHeartbeat() {
	if !m.heartbeat {
		return
	}

	today := time.Now().Format("2006-01-02")

	if _, err := os.Stat(heartbeatFile); errors.Is(err, os.ErrNotExist) {
		m.updateHeartbeat(today)
		return
	}

	bytes, err := os.ReadFile(heartbeatFile)
	if err != nil {
		m.pusher.Push(fmt.Sprintf("%s: FAILED", m.name), fmt.Sprintf("Failed to read heartbeat file, err = %v", err))
		return
	}
	if string(bytes) != today {
		m.updateHeartbeat(today)
	}
}

func (m *Monitor) updateHeartbeat(today string) {
	err := os.WriteFile(heartbeatFile, []byte(today), 0644)
	if err != nil {
		m.pusher.Push(
			fmt.Sprintf("%s: FAILED", m.name),
			fmt.Sprintf("Failed to write to heartbeat file, err = %v", err))
		return
	}
	m.pusher.Push(
		fmt.Sprintf("%s: HEARTBEAT", m.name),
		fmt.Sprintf("I'm still alive!\n%s", strings.Join(m.getNetInfo(), "\n")))
}

func (m *Monitor) Check() {
	m.doHeartbeat()

	netInfo := m.getNetInfo()
	netInfoStr := strings.Join(netInfo, "\n")

	// tmpFile not exists
	if _, err := os.Stat(netInfoFile); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(netInfoFile, []byte(strings.Join(netInfo, "\n")), 0644)
		if err != nil {
			m.pusher.Push(fmt.Sprintf("%s: FAILED", m.name), fmt.Sprintf("Failed to write to network info file, err = %v", err))
			return
		}
		m.pusher.Push(fmt.Sprintf("%s: INIT", m.name), strings.Join(netInfo, "\n"))
		return
	}

	bytes, err := os.ReadFile(netInfoFile)
	if err != nil {
		m.pusher.Push(fmt.Sprintf("%s: FAILED", m.name), fmt.Sprintf("Failed to read network info file, err = %v", err))
		return
	}

	if string(bytes) != netInfoStr {
		m.updateNetInfo(netInfoStr)
	}

}

func (m *Monitor) updateNetInfo(netInfoStr string) {
	err := os.WriteFile(netInfoFile, []byte(netInfoStr), 0644)
	if err != nil {
		m.pusher.Push(fmt.Sprintf("%s: FAILED", m.name), fmt.Sprintf("Failed to write to network info file, err = %v", err))
		return
	}
	m.pusher.Push(fmt.Sprintf("%s: IPs CHANGED", m.name), netInfoStr)
}

func (m *Monitor) getNetInfo() []string {
	var stringArr []string
	links, err := netlink.LinkList()
	if err != nil {
		m.pusher.Push(fmt.Sprintf("%s: FAILED", m.name), fmt.Sprintf("Failed to list link, err = %v", err))
		return stringArr
	}
	for _, link := range links {
		linkName := link.Attrs().Name
		if len(m.devices) > 0 && !stringArrContains(m.devices, linkName) {
			continue
		}
		addrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
		if err != nil {
			m.pusher.Push(fmt.Sprintf("%s: FAILED", m.name), fmt.Sprintf("Failed to list addr for link %s, err = %v", linkName, err))
			continue
		}
		addrsArr := make([]string, 0, len(addrs))
		for _, addr := range addrs {
			addrsArr = append(addrsArr, addr.IP.String())
		}
		sort.Strings(addrsArr)
		stringArr = append(stringArr, fmt.Sprintf("- %s: %s", linkName, strings.Join(addrsArr, ", ")))
	}
	return stringArr
}

func stringArrContains(arr []string, target string) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}
