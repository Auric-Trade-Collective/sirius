//go:build linux
package main

import (
	"log/slog"
	"os"

	"github.com/YendisFish/sirius/alpha/config"
	"github.com/YendisFish/sirius/alpha/monitor"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/sys/unix"
)

func main() {
	slog.Info("Initializing alpha")

	err := mountFs()
	if err != nil {
		return
	}

	slog.Info("Reading configuration...")

	fle, err := os.ReadFile("/etc/alpha/alpha.toml")
	if err != nil {
		slog.Error("Problems reading alpha.toml: " + err.Error())
		return
	}

	var config config.Config
	err = toml.Unmarshal(fle, &config)
	if err != nil {
		slog.Error("Problems parsing alpha.toml: " + err.Error())
		return
	}

	slog.Info("Finished initializing alpha")
	slog.Info("Starting system...")

	mon := monitor.NewMonitor()
	mon.Config = config

	for n, entry := range config.Host {
		go mon.CreateDaemonProcess(n, entry)
	}

	for {
		mon.RunCycle()
	}
}

func mountFs() error {
	err := unix.Mount("none", "/dev", "devtmpfs", unix.MS_NOSUID, "")
	if err != nil {
		slog.Error("Could not mount filesystem...")
		return err
	}

	err = unix.Mount("proc", "/proc", "proc", 0, "")
    if err != nil {
        return err
    }

    err = unix.Mount("sysfs", "/sys", "sysfs", 0, "")
    if err != nil {
        return err
    }

	return nil
}
