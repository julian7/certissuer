package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

var (
	version = "SNAPSHOT"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("cannot create logger")
	}

	app := NewApp(logger)
	if err := app.Command().Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "\nError: %v\n", err.Error())
		os.Exit(1)
	}
}
