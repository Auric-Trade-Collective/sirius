package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	// "golang.org/x/term"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

var cmd = flag.String("c", "", "Use to execute a command once")

func main() {
	flag.Parse()

	runner, err := interp.New(interp.Interactive(true), interp.StdIO(os.Stdin, os.Stdout, os.Stderr))
	if err != nil {
		slog.Info("Could not initialize runner Reason: " + err.Error())
		fmt.Println("Could not initialize runner Reason: " + err.Error())
	}

	parser := syntax.NewParser()

	if *cmd == "" {
		runInteractive(runner, parser)
	}
}

func runInteractive(runner *interp.Runner, parser *syntax.Parser) {
	fmt.Print("$ ")
	for stmts, err := range parser.InteractiveSeq(os.Stdin) {
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if parser.Incomplete() {
			fmt.Print("> ")
			continue
		}

		ctx := context.Background()
		for _, stmt := range stmts {
			err := runner.Run(ctx, stmt)
			if err != nil {
				fmt.Println(err)
			}
			if runner.Exited() {
				fmt.Println()
				os.Exit(0)
			}
		}

		fmt.Print("$ ")
	}
}
