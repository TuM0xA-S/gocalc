package main

import (
	"flag"
	"os"

	"github.com/TuM0xA-S/gocalc"
)

func main() {
	script := flag.Bool("s", false, "script mode")
	precision := flag.Int("p", 2, "precision")
	flag.Parse()

	gocalc.NewInterpreter(!*script, *precision).Start(os.Stdin, os.Stdout)
}
