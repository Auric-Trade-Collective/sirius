package monitor

import (
	"log/slog"
	"os"
)


func CreateOrGetJournal(name string) (*os.File, error) {
	if fs, err := os.Open("/var/log/alpha/" + name + ".journal"); err == nil {
		return fs, nil
	}

	fs, err := os.Create("/var/log/alpha/" + name + ".journal")
	if err != nil {
		slog.Error("Could not get or create journal for: " + name)
		return nil, err
	}

	return fs, nil
}
