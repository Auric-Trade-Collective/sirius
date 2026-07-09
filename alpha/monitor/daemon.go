package monitor

import (
	"os/exec"

	"github.com/YendisFish/sirius/alpha/config"
)

type AlphaDaemon struct {
	Info *config.Entry
	Cmd *exec.Cmd
	Send chan func()
	Recv chan ProcessEvent
	Name string
}

func (a AlphaDaemon) GetType() ServiceType { return ServiceTypeDaemon }
