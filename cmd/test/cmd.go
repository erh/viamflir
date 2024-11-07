package main

import (
	"context"
	"flag"
	"fmt"

	"go.viam.com/rdk/logging"

	"github.com/erh/viamflir"
)

func main() {
	err := realMain()
	if err != nil {
		panic(err)
	}
}

func realMain() error {
	logger := logging.NewLogger("viamflir-test")

	debug := false

	flag.BoolVar(&debug, "debug", debug, "")

	if debug {
		logger.SetLevel(logging.DEBUG)
	}

	ip, err := viamflir.FindIP(context.Background(), logger)
	if err != nil {
		return err
	}
	fmt.Printf("found ip: %v\n", ip)
	return nil
}
