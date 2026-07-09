package reader

import (
	"fmt"
	"log/slog"
	"os"

	"golang.org/x/term"
)

func ReadPassword(header string) string {
	fmt.Print(header)

	fd := int(os.Stdin.Fd())
	bytePassword, err := term.ReadPassword(fd)
	if err != nil {
		slog.Error("Could not read password: " + err.Error())
		return ""
	}

	fmt.Println()

	return string(bytePassword)
}
