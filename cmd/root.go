package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
)

func Die(err error) int {
	log.Fatal("Error: ", err.Error())

	return 1
}

func Main(output io.Writer) int {
	conf, err := InitConfig(output)
	if err != nil {
		return Die(err)
	}

	if conf.Examples {
		fmt.Fprintln(output, Examples)
		os.Exit(0)
	}

	tp := NewTP(conf)

	if err := tp.ProcessTimestamps(); err != nil {
		return Die(err)
	}

	return 0
}
