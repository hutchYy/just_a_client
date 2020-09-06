package main

import (
	"net"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type Session struct {
	id string
	//id int64
	wg sync.WaitGroup
	//deviceName string
	ClientIdentifier string
	conn             net.Conn
	closed           bool
}

var globalSessionID int64
var globalSessionCount int32

var sessions []*Session

func RemoveIndex(s []*Session, index int) []*Session {
	return append(s[:index], s[index+1:]...)
}

func NewSession() *Session {
	var session Session
	id, _ := uuid.NewV4()
	session.id = id.String()
	atomic.AddInt32(&globalSessionCount, 1)
	return &session
}

func startListener(host string, port int) {
	// Listen for incoming connections.
	addr := host + ":" + strconv.Itoa(port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("Error listening "+addr, err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	log.Println("Listening on " + addr)
	go func(l net.Listener) {
		for {
			c, err := l.Accept()
			if err != nil {
				log.Printf("Closing listener goroutine")
			}
			newConns <- c
		}
	}(l)

	handleClient(l)
}

func handleClient(l net.Listener) {
	for {
		select {
		case v, ok := <-stop:
			if v == true && ok == true {
				log.Printf("Clossing opened sessions")
				for i, session := range sessions {
					session.conn.Close()
					sessions = RemoveIndex(sessions, i)
				}
				log.Printf("Stopping listener")
				l.Close()
			}
			return
		case c := <-newConns:
			log.Println("New client connected")
			session := NewSession()
			sessions = append(sessions, session)
			go session.startShell(c)
		}
	}
}
