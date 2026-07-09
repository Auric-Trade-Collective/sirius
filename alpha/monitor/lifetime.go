package monitor

import (
	"log/slog"
	"os"
	"os/exec"

	"github.com/YendisFish/sirius/alpha/config"
	"github.com/YendisFish/sirius/alpha/kernel"
)

func (m *Monitor) RunCycle() {
	for _, ac := range m.ByPid {
		switch ac := ac.(type) {
			case AlphaDaemon:
				handleDaemonEvent(ac)
			default:
				continue
		}
	}
}

func handleDaemonEvent(service AlphaDaemon) {
	select {
		case procE := <- service.Recv:
			switch event := procE.(type) {
			case ProcessExit:
				slog.Info("Process exited: " + service.Name)
				_ = event
			}
		default:
	}
}

func daemonMonitor(me AlphaDaemon) {
	exit := make(chan error, 1)
	go func() {
		exit <- me.Cmd.Wait()
	}()

	for {
		select {
			case action := <- me.Send:
				action()
			case <- exit:
				me.Recv <- ProcessExit{
					Code: me.Cmd.ProcessState.ExitCode(),
				}
			default:
				continue
		}
	}
}

func handleDepdendencies(entry config.Entry) {
	for _, dev := range entry.NeedsDev {
		if dev == "/dev/console" {
			entry.NeedsTTY = new(bool)
			*entry.NeedsTTY = true
		}
		kernel.WaitForDevice(dev)
	}
}

func handleIO(cmd *exec.Cmd, name string, entry config.Entry) error {
	if entry.NeedsTTY != nil && *entry.NeedsTTY == true {
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

	return nil
}
