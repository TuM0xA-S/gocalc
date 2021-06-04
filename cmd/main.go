package main

import (
	"os"

	"github.com/TuM0xA-S/gocalc"
)

func main() {
	gocalc.NewInterpreter(true, 2).Start(os.Stdin, os.Stdout)
}
