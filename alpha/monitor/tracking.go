package monitor

import (
	"slices"
	"syscall"
)

func (m *Monitor) killDependenciesOf(name string) {
	m.mutex.Lock()

	for _, proc := range m.ByName {
		switch proc.GetType() {
			case ServiceTypeDaemon:
				daemon := proc.(AlphaDaemon)

				if daemon.Info.NeedsDep != nil && slices.Contains(daemon.Info.NeedsDep, name) {
					daemon.Cmd.Process.Signal(syscall.SIGTERM)
				}
		}
	}

	m.mutex.Unlock()
}
