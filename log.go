package main

import "log"

type logger struct {
	level int
}

var l logger

func initLogger() {
	switch config.logLevel {
	case "all":
		l = logger{0}
	case "info":
		l = logger{1}
	case "error":
		l = logger{2}
	case "fatal":
		l = logger{3}
	case "none":
		l = logger{4}
	default:
		l = logger{0}
	}
}

func (l logger) Infof(format string, args ...interface{}) {
	if l.level <= 1 {
		log.SetPrefix("INFO: ")
		log.Printf(format, args...)
	}
}

func (l logger) Errorf(format string, args ...interface{}) {
	if l.level <= 2 {
		log.SetPrefix("ERROR: ")
		log.Printf(format, args...)
	}
}

func (l logger) Fatalf(format string, args ...interface{}) {
	if l.level <= 3 {
		log.SetPrefix("FATAL: ")
		log.Fatalf(format, args...)
	}
}
