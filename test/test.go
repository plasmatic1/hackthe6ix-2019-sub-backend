package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const ADDR = "0.0.0.0:80"

func waitForLines(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			fmt.Printf("Reader error: %s\n", err.Error())
		} else {
			fmt.Printf("Received line: \"%s\"\n", line)
		}
	}
}

func main() {
	conn, err := net.Dial("tcp", ADDR)
	if err != nil {
		fmt.Printf("Listen Error: %s\n", err.Error())
	} else {
		fmt.Printf("Opened test client at addr %s\n", ADDR)

		go waitForLines(conn)

		reader := bufio.NewReader(os.Stdin)
		for {
			line, _, _ := reader.ReadLine()
			line = append(line, '\n')
			if n, err := conn.Write(line); err != nil {
				fmt.Printf("Send to server error: %s\n", err.Error())
			} else {
				fmt.Printf("Wrote %d bytes\n", n)
			}
		}
	}
}
