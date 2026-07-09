//go:build linux

package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/YendisFish/sirius/guarddog/reader"
	"github.com/YendisFish/sirius/toybox/types"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/sys/unix"
)

func main() {
	waitForTTY()

	console, err := os.OpenFile("/dev/console", os.O_RDWR, 0)
	if err != nil {
		slog.Error("Error initializing TTY")
	}

	fd := int(console.Fd())
	unix.Dup2(fd, 0)
	unix.Dup2(fd, 1)
	unix.Dup2(fd, 2)

	fmt.Print("username: ")
	scanner := bufio.NewReader(os.Stdin)

	uname, _, err := scanner.ReadLine()
	if err != nil {
		slog.Error("Could not read standard input, Reason: " + err.Error())
	}

	passwd := reader.ReadPassword("password: ")
	login(string(uname), passwd)
}

func login(username string, password string) {
	usrs, err := types.ReadPasswd()
	if err != nil {
		slog.Error("Couldn't read password!")
		panic("Password failure")
	}

	for _, usr := range usrs {
		if usr.Username == username {
			if len(password) > 0 {
				//perform hashing and compare
				break
			}

			//we get a shell now!
			fmt.Println("Success, starting shell!")
		}
	}
}

func waitForTTY() {
	if _, err := os.Stat("/dev/console"); err == nil {
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error(err.Error())
	}
	defer watcher.Close()

	err = watcher.Add("/dev/console")
	if err != nil {
		slog.Error(err.Error())
	}

	for {
		select {
			case event, ok := <- watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Create == fsnotify.Create && event.Name == "/dev/console" {
					return
				}
			case error, ok := <-watcher.Errors:
				if !ok {
					slog.Error("ERROR")
				}

				if error != nil {
					slog.Error("ERROR: " + error.Error())
				}
		}
	}
}
