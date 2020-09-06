package main

import (
	"flag"
	"fmt"
	"os"
)

var UtiliyFlagSet *flag.FlagSet
var ListenerFlagSet *flag.FlagSet

var host string
var port int
var logging bool
var verbosity string

var defaultMessage = `
Usage : just_a_server [host] port [logging] [verbosity]
[] = optional
Commands:
  host	  		Set listening host (0.0.0.0 by default)
  port	  		Set listening port (1111 by default)
  logging	  	Enable logging
  verbosity	  	Set verbosity
Run 'just_a_server -h' for more information on a command.`

func initFlags() {
	ListenerFlagSet = flag.NewFlagSet("listener", flag.ExitOnError)
	ListenerFlagSet.StringVar(&host, "host", "0.0.0.0", "Set listening ip")
	ListenerFlagSet.IntVar(&port, "port", 1111, "Set listening port")
	UtiliyFlagSet = flag.NewFlagSet("", flag.ExitOnError)
	UtiliyFlagSet.BoolVar(&logging, "logging", false, "Enable logging")
	UtiliyFlagSet.StringVar(&verbosity, "verbosity", "info", "Set log verbosity several allowed [panic, fatal, error, warning, info, debug, trace]")
}

func argParser() {
	initFlags()
	if len(os.Args) < 2 {
		switch os.Args[1] {
		case "listener":
			ListenerFlagSet.Parse(os.Args[2:])
		default:
			fmt.Println(defaultMessage)
			os.Exit(1)
		}
	}
}
