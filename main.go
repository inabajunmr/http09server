package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func main() {
	service := ":80"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		handle(listener)
	}
}

func handle(listener *net.TCPListener) {
	conn, err := listener.Accept()
	if err != nil {
		return
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	line = strings.Trim(line, "\n")
	checkError(err)
	req := strings.Split(line, " ")
	wd, err := os.Getwd()
	checkError(err)
	if len(req) < 2 {
		response(wd+"/"+"badrequest.html", conn)
		return
	}
	if req[0] != "GET" {
		response(wd+"/"+"badrequest.html", conn)
		return
	}

	// '/' to index.html
	if len(req[1][strings.Index(req[1], "/")+1:]) == 0 {
		response(wd+"/"+"index.html", conn)
		return
	}

	path := wd + "/" + req[1][strings.Index(req[1], "/")+1:]
	f, err := os.Open(path) // directory traversal
	defer f.Close()
	if err != nil {
		response(wd+"/"+"notfound.html", conn)
		return
	}

	fInfo, err := os.Stat(path)
	checkError(err)
	if fInfo.IsDir() {
		response(wd+"/"+"notfound.html", conn)
		return
	}

	response(path, conn)
}

func response(path string, conn net.Conn) {
	f, err := os.Open(path)
	checkError(err)
	b, err := ioutil.ReadAll(f)
	checkError(err)
	f.Close()
	conn.Write([]byte(b))
	conn.Close()

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}
