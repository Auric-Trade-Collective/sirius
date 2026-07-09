//go:build linux

package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/YendisFish/sirius/guarddog/reader"
	"github.com/YendisFish/sirius/toybox/types"
)

func main() {
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
