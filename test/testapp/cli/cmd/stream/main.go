package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	socketPath := os.Args[1]
	_ = os.Remove(socketPath)
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "stream: listen %s: %v\n", socketPath, err)
		os.Exit(1)
	}
	if err := os.Chmod(socketPath, 0666); err != nil {
		fmt.Fprintf(os.Stderr, "stream: chmod %s: %v\n", socketPath, err)
		os.Exit(1)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "stream: accept: %v\n", err)
			continue
		}
		go serve(conn)
	}
}

func serve(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	header, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	fields := strings.Fields(header)
	if !validProxyHeader(fields) {
		fmt.Fprintf(conn, "error invalid proxy header: %q\n", header)
		return
	}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		fmt.Fprintf(conn, "stream %s %s\n", fields[2], strings.TrimSpace(line))
	}
}

func validProxyHeader(fields []string) bool {
	if len(fields) != 6 || fields[0] != "PROXY" {
		return false
	}
	v4 := fields[1] == "TCP4"
	if !v4 && fields[1] != "TCP6" {
		return false
	}
	for _, field := range fields[2:4] {
		ip := net.ParseIP(field)
		if ip == nil || (ip.To4() != nil) != v4 {
			return false
		}
	}
	for _, field := range fields[4:6] {
		port, err := strconv.Atoi(field)
		if err != nil || port < 0 || port > 65535 {
			return false
		}
	}
	return true
}
