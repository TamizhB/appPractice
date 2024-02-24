package messageHandlers

import (
	"fmt"
	"net"
)

const (
	port = 8888
)

// Serve
func Serve() {
	listenAddr, err := net.ResolveUDPAddr("udp", ":8888")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	// Create a UDP listener
	conn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}

	// defer conn.Close() ---> runs indefinitely

	// Buffer to receive data
	buffer := make([]byte, 1024)

	// Infinite loop to receive data
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading data:", err)
			continue
		}

		fmt.Printf("Received %d bytes from %s:", n, addr)
		fmt.Println(string(buffer[:n]))
	}
}
