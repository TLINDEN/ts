package main

import (
	"os"

	"github.com/tlinden/ts/cmd"
)

func main() {
	os.Exit(cmd.Main(os.Stdout))
}
