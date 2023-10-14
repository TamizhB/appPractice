package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

func main() {
	serverAddress := "localhost"
	serverPort := 12345
	message := "getconfig"

	// Resolve server address
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", serverAddress, serverPort))
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	// Create a UDP connection for the client
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error creating UDP connection:", err)
		return
	}
	defer conn.Close()

	// Convert the message to bytes
	messageBytes := []byte(message)

	// Send the message to the server
	_, err = conn.Write(messageBytes)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	fmt.Printf("Sent to %s:%d: %q\n", serverAddress, serverPort, message)

	// Create a buffer to receive the response
	responseBuf := make([]byte, 16)

	// Set a deadline for waiting for a response
	//conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Read the response from the server
	n, err := conn.Read(responseBuf)
	if err != nil {
		fmt.Println("Error receiving response:", err)
		return
	}

	fmt.Printf("Received response: %x\n", responseBuf[:n])

	// Parse the binary buffer received from the server
	angle := int32(binary.LittleEndian.Uint32(responseBuf[0:4]))
	rotate := responseBuf[4] == 1
	frequency := string(responseBuf[5:n])

	fmt.Println("Parsed Data:")
	fmt.Printf("angle: %d\n", angle)
	fmt.Printf("rotate: %v\n", rotate)
	fmt.Printf("frequency: %s\n", frequency)

	go func() {

		buf := make([]byte, 1024)
		for {
			n, addr, err := conn.ReadFrom(buf)
			if err != nil {
				fmt.Printf("Error receiving notification: %v\n", err)
				continue
			}

			notification := string(buf[:n])
			fmt.Printf("Received notification from %s: %s\n", addr.String(), notification)
		}
	}()

	// Keep the subscriber running
	select {}
}
