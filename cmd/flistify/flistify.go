package main

import (
	"fmt"
	"os"

	"github.com/omarabdul3ziz/flistify/internal/builder"
)

const (
	BUILD = "build"
	PUSH  = "push"
	RUN   = "run"
)

func main() {
	// TODO: use logger
	// TODO: use cli tool

	if len(os.Args) < 2 {
		fmt.Printf("missing args, got: %v\n", os.Args)
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case BUILD:
		err = builder.Build(os.Args[2:])
	case PUSH:
		// err = push(os.Args[2:])
	case RUN:
		// err = run(os.Args[2:])
	default:
		err = fmt.Errorf("\"%v\" is not supported", os.Args[1])
	}

	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
		os.Exit(1)
	}
}
