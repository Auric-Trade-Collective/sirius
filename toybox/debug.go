package toybox

import (
	"fmt"
	"log/slog"
	"os"
)

func ReadDir(dir string) {
	fles, err := os.ReadDir(dir)

	if err != nil {
		slog.Error("Could not list directory for debug info Reason: " + err.Error())
	}

	for _, fle := range fles {
		fmt.Println(fle.Name())
	}
}
