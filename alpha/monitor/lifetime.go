package monitor

import (
	"log/slog"
	"os"
	"os/exec"
	"syscall"

	"github.com/YendisFish/sirius/alpha/config"
	"github.com/YendisFish/sirius/alpha/kernel"
)

func (m *Monitor) RunCycle() {
	m.mutex.Lock()
	for _, ac := range m.ByPid {
		switch ac := ac.(type) {
			case AlphaDaemon:
				m.handleDaemonEvent(ac)
			default:
				continue
		}
	}
	m.mutex.Unlock()
}

func (m *Monitor) handleDaemonEvent(service AlphaDaemon) {
	select {
		case procE := <- service.Recv:
			switch event := procE.(type) {
			case ProcessExit:
				go m.killDependenciesOf(service.Name)
				slog.Info("Process exited: " + service.Name)
				go m.handleTeardown(service)
				_ = event
			case ProcessCrash:
				go m.killDependenciesOf(service.Name)
				go m.handleTeardown(service)
			}
		default:
	}
}

func daemonMonitor(me AlphaDaemon) {
	me.Cmd.Wait()

	me.Recv <- ProcessExit{
		Code: me.Cmd.ProcessState.ExitCode(),
	}
}

func (m *Monitor) handleDepdendencies(entry *config.Entry) error {
	if entry.NeedsDep != nil {
		for _, proc := range entry.NeedsDep {
			if !m.isDaemonRunning(proc) {
				entry, err := m.Config.FindEntryByName(proc)
				if err != nil {
					return err
				}

				m.CreateDaemonProcess(proc, entry)
			}
		}
	}

	if entry.NeedsDev != nil {
		for _, dev := range entry.NeedsDev {
			if dev == "/dev/console" {
				entry.NeedsTTY = new(bool)
				*entry.NeedsTTY = true
			}
			kernel.WaitForDevice(dev)
		}
	}

	return nil
}

func handleIO(cmd *exec.Cmd, name string, entry *config.Entry) error {
	if entry.NeedsTTY != nil && *entry.NeedsTTY != true {
		journal, err := CreateOrGetJournal(name)
		if err != nil {
			slog.Error("Could not create journal for: " + entry.Name)
			return err
		}

		cmd.Stdout = journal
		cmd.Stderr = journal
		cmd.Stdin = nil

		return nil
	}

	console, err := os.OpenFile("/dev/console", os.O_RDWR, 0)
	if err != nil {
		slog.Error("Error initializing TTY")
		return err
	}

	cmd.Stdin = console
	cmd.Stdout = console
	cmd.Stderr = console

	cmd.SysProcAttr = &syscall.SysProcAttr{
	    Setsid: true,
	    Setctty: true,
	    Ctty: 0,
	}

	return nil
}

func (m *Monitor) handleTeardown(s Service) {
	m.mutex.Lock()

	switch s.GetType() {
		case ServiceTypeDaemon:
			daemon := s.(AlphaDaemon)
			delete(m.ByName, daemon.Name)
			delete(m.ByPid, daemon.Cmd.Process.Pid)

			if daemon.Info.OnExit == nil || *daemon.Info.OnExit == "restart" {
				slog.Info("Restarting service: " + daemon.Name)
				go m.CreateDaemonProcess(daemon.Name, daemon.Info)
			}
	}

	m.mutex.Unlock()
}
