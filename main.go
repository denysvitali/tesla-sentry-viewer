package main

import (
	"github.com/alexflint/go-arg"
	sentry "github.com/denysvitali/tesla-sentry-viewer/pkg"
	"github.com/sirupsen/logrus"
)

var args struct {
	Debug          *bool  `arg:"-D,--debug"`
	EventDirectory string `arg:"positional,required" help:"A directory containing an event.json file"`
}

func main() {
	arg.MustParse(&args)
	logger := logrus.New()
	if args.Debug != nil && *args.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	err := sentry.ProcessVideo(args.EventDirectory)
	if err != nil {
		logger.Fatalf("unable to parse: %v", err)
	}
}
