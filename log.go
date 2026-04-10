package main

import "log"

type logger struct {
	level int
}

func initLogger() logger {
	switch config.logLevel {
	case "all":
		return logger{0}
	case "info":
		return logger{1}
	case "error":
		return logger{2}
	case "fatal":
		return logger{3}
	case "none":
		return logger{4}
	default:
		return logger{4}
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
