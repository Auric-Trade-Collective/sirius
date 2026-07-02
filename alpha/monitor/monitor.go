package monitor

import (
	"log/slog"
	"os/exec"
	"strconv"
	"sync"

	"github.com/YendisFish/sirius/alpha/config"
)

type AlphaProcess struct {
	Info config.Entry
	Cmd *exec.Cmd
	Send chan func()
	Recv chan ProcessEvent
}

type ProcessEvent interface {}

type ProcessExit struct {
	Code int
}

type Monitor struct {
	mutex sync.RWMutex
	Host map[int]*AlphaProcess
}

func NewMonitor() *Monitor {
	ret := &Monitor{
		Host: make(map[int]*AlphaProcess),
	}

	return ret
}

func (m *Monitor) CreateHostProcess(name string, entry config.Entry) {
	cmd := exec.Command(entry.Name, entry.Args)

	journal, err := CreateOrGetJournal(name)
	if err != nil {
		slog.Error("Could not create journal for: " + entry.Name)
		return
	}

	cmd.Stdout = journal
	cmd.Stderr = journal
	cmd.Stdin = nil

	err = cmd.Start()
	if err != nil {
		slog.Error("Couldn't start process: " + entry.Name + " Reason: " + err.Error())
		return
	}

	m.mutex.Lock()
	if _, exists := m.Host[cmd.Process.Pid]; !exists {
		alp := &AlphaProcess{
			Info: entry,
			Cmd: cmd,
			Send: make(chan func()),
			Recv: make(chan ProcessEvent),
		}

		m.Host[cmd.Process.Pid] = alp

		go processMonitor(alp)
	} else {
		slog.Info("Ambiguous PID found: " + strconv.Itoa(cmd.Process.Pid))
		return
	}
}

func processMonitor(me *AlphaProcess) {
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
