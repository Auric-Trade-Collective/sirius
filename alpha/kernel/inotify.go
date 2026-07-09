package kernel

import (
	"log/slog"
	"os"

	"github.com/fsnotify/fsnotify"
)

func WaitForDevice(name string) {
	if _, err := os.Stat(name); err == nil {
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error(err.Error())
	}
	defer watcher.Close()

	err = watcher.Add(name)
	if err != nil {
		slog.Error(err.Error())
	}

	for {
		select {
			case event, ok := <- watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Create == fsnotify.Create && event.Name == name {
					return
				}
			case error, ok := <-watcher.Errors:
			//this entire case needs to be fixed
				if !ok {
					slog.Error("ERROR")
				}

				if error != nil {
					slog.Error("ERROR: " + error.Error())
				}
		}
	}
}
