package pkg

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/vishvananda/netlink"
)

const tmpFile = "/tmp/ip-monitor.txt"

type Monitor struct {
	name string

	pusher *ServerChan
}

func NewMonitor(name string, pusher *ServerChan) *Monitor {
	return &Monitor{
		name:   name,
		pusher: pusher,
	}
}

func (m *Monitor) Check() {
	links, err := netlink.LinkList()
	if err != nil {
		m.pusher.Push(fmt.Sprintf("%s: FAILED", m.name), fmt.Sprintf("Failed to list link, err = %v", err))
		return
	}

	sb := new(strings.Builder)
	for _, link := range links {
		linkName := link.Attrs().Name
		sb.WriteString(fmt.Sprintf("- %s: ", linkName))
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
		sb.WriteString(strings.Join(addrsArr, ", "))
		sb.WriteByte('\n')
	}

	// tmpFile not exists
	if _, err := os.Stat(tmpFile); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(tmpFile, []byte(sb.String()), 0644)
		if err != nil {
			m.pusher.Push(fmt.Sprintf("%s: FAILED", m.name), fmt.Sprintf("Failed to write to tmp file, err = %v", err))
			return
		}
		m.pusher.Push(fmt.Sprintf("%s: INIT", m.name), sb.String())
		return
	}

	bytes, err := os.ReadFile(tmpFile)
	if err != nil {
		m.pusher.Push(fmt.Sprintf("%s: FAILED", m.name), fmt.Sprintf("Failed to read tmp file, err = %v", err))
		return
	}
	if string(bytes) != sb.String() {
		err := os.WriteFile(tmpFile, []byte(sb.String()), 0644)
		if err != nil {
			m.pusher.Push(fmt.Sprintf("%s: FAILED", m.name), fmt.Sprintf("Failed to write to tmp file, err = %v", err))
			return
		}
		m.pusher.Push(fmt.Sprintf("%s: IPs CHANGED", m.name), sb.String())
	}
}
