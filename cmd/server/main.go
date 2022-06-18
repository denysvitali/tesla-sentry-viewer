package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/denysvitali/tesla-sentry-viewer/pkg/server"
	"github.com/sirupsen/logrus"
)

var args struct {
	Debug     *bool  `arg:"--debug,-D"`
	Directory string `arg:"positional,required"`
	Listen    string `arg:"-l,--listen" default:"127.0.0.1:8150"`
	Version   *bool  `arg:"-v,--version"`
}

const SoftwareName = "tesla-sentry-viewer"

var Version = "dev"

func main() {
	arg.MustParse(&args)

	if args.Version != nil && *args.Version {
		fmt.Printf("%s %s", SoftwareName, Version)
	}

	logger := logrus.New()

	if args.Debug != nil && *args.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	s, err := server.New(args.Directory)
	if err != nil {
		logger.Fatalf("unable to create server: %v", err)
	}

	s.SetLogger(logger)
	err = s.Listen(args.Listen)
	if err != nil {
		logger.Fatalf("unable to start server: %v", err)
	}
}
