//go:build linux

package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"syscall"

	toybox_reader "github.com/YendisFish/sirius/guarddog/reader"
	toybox_types "github.com/YendisFish/sirius/toybox/types"
)

func main() {
	fmt.Print("username: ")
	scanner := bufio.NewReader(os.Stdin)

	uname, _, err := scanner.ReadLine()
	if err != nil {
		slog.Error("Could not read standard input, Reason: " + err.Error())
	}

	passwd := toybox_reader.ReadPassword("password: ")
	login(string(uname), passwd)
}

func login(username string, password string) {
	usrs, err := toybox_types.ReadPasswd()
	if err != nil {
		slog.Error("Couldn't read password! Reason: " + err.Error())
		os.Exit(1)
	}

	for _, usr := range usrs {
		if usr.Username == username {
			if len(password) > 0 {
				//perform hashing and compare
				break
			}

			//we get a shell now!
			fmt.Println("Success, starting shell!")
			transfer(usr)
		}
	}
}

func transfer(user toybox_types.PasswdUser) {
	//eventually groups will need to be set as well

	if err := syscall.Setresgid(user.GID, user.GID, user.GID); err != nil {
		panic("Couldn't transfer program ownership to user")
	}

	if err := syscall.Setresuid(user.UID, user.UID, user.UID); err != nil {
		panic("Couldn't transfer program ownership to user")
	}

	syscall.Exec(user.Shell, []string{user.Shell}, os.Environ())
}
