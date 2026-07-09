package monitor

import (
	"log/slog"
	"os/exec"
	"strconv"
	"sync"

	"github.com/YendisFish/sirius/alpha/config"
)

type ServiceType int
const (
	ServiceTypeProcess ServiceType = iota
	ServiceTypeContainer
	ServiceTypeDaemon
)

type Service interface {
	GetType() ServiceType
}

type Monitor struct {
	mutex sync.RWMutex
	ByPid map[int]Service
	ByName map[string]Service
}

func NewMonitor() *Monitor {
	ret := &Monitor{
		ByPid: make(map[int]Service),
		ByName: make(map[string]Service),
	}

	return ret
}

func (m *Monitor) CreateDaemonProcess(name string, entry config.Entry) {
	handleDepdendencies(entry) //blocks until necessary processes and devices are hooked up
	cmd := exec.Command(entry.Name, entry.Args)

	err := handleIO(cmd, name, entry)
	if err != nil {
		slog.Error("Could not successfully start: " + name)
		return
	}

	err = cmd.Start()
	if err != nil {
		slog.Error("Couldn't start process: " + entry.Name + " Reason: " + err.Error())
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, exists := m.ByName[name]; !exists {
		alp := AlphaDaemon{
			Info: entry,
			Cmd: cmd,
			Send: make(chan func()),
			Recv: make(chan ProcessEvent),
			Name: name,
		}

		m.ByPid[cmd.Process.Pid] = alp
		m.ByName[name] = alp

		go daemonMonitor(alp)
	} else {
		slog.Info("Ambiguous Name found: " + strconv.Itoa(cmd.Process.Pid))
		return
	}
}
