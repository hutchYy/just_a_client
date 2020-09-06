package main

import (
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

var start chan bool
var stop chan bool
var newConns chan net.Conn

func main() {
	stop = make(chan bool)
	start = make(chan bool)
	newConns = make(chan net.Conn)
	argParser()
	if logging {
		oldStd = os.Stdout
		stdr, stdw, _ = os.Pipe()
		os.Stdout = stdw
		initializeLogging("file.log")
		exist := stringInSlice(verbosity, []string{"panic", "fatal", "error", "warning", "info", "debug", "trace"})
		if exist {
			lvl, _ := log.ParseLevel(verbosity)
			log.SetLevel(lvl)
			log.Println("Log Level Set To : " + verbosity)
		} else {
			fmt.Fprintln(oldStd, "Only the following levels are allowed : \n\n\tpanic\n\tfatal\n\terror\n\twarning\n\tinfo\n\tdebug\n\ttrace")
			os.Exit(1)
		}
	}
	startListener(host, port)
}
