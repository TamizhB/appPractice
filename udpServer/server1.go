package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {
	serverPort := 23456

	// Create a UDP connection
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println("Error creating UDP connection:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("UDP Server listening on %s\n", conn.LocalAddr())

	// Create a buffer to receive messages
	buf := make([]byte, 16)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		fmt.Printf("Received from %s(%d bytes): %s\n", remoteAddr, n, string(buf))

		if string(buf[:n]) == "getconfig" {
			// Construct a binary buffer with integer, boolean, and string values
			responseBuf := make([]byte, 16)

			// angle
			binary.LittleEndian.PutUint32(responseBuf[0:], 42)

			// rotate
			responseBuf[4] = 1

			// carrier frequency
			copy(responseBuf[5:], []byte("70GHz"))

			// Send the binary buffer back to the client
			_, err = conn.WriteToUDP(responseBuf, remoteAddr)
			if err != nil {
				fmt.Println("Error sending message:", err)
			}
		}
		if bytes.Equal(buf, []byte("subscribe")) {
			go notifyAsync(conn, remoteAddr)
		}
	}

}

func notifyAsync(conn *net.UDPConn, remoteAddr *net.UDPAddr) {
	for {
		notification := "notification at " + time.Now().String()
		_, err := conn.WriteToUDP([]byte(notification), remoteAddr)
		if err != nil {
			fmt.Printf("Error sending notification: %v\n", err)
		} else {
			fmt.Println("Sent " + notification)
		}
		time.Sleep(time.Minute)
	}
}
