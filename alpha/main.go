//go:build linux

package main

import (
	"log/slog"
	"os"

	"github.com/YendisFish/sirius/alpha/config"
	"github.com/YendisFish/sirius/alpha/kernel"
	"github.com/YendisFish/sirius/alpha/monitor"
	"github.com/pelletier/go-toml/v2"
)

func main() {
	slog.Info("Initializing alpha")

	if err := kernel.MountFs(); err != nil {
		slog.Info("Could not initialize ramfs")
		return
	}

	if err := kernel.MountZFS(); err != nil {
		slog.Info("Could not initialize zfs")
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
