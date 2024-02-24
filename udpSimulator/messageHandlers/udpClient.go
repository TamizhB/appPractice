package messageHandlers

import (
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
)

// Send
func sendMessage(ctx *gin.Context) {
	// Replace with the receiver's IP address and port
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	// Create a connection to the server
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	defer conn.Close()

	// Data to send
	data := []byte("Hello, world!")

	// Send the data
	n, err := conn.Write(data)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}
	fmt.Printf("Sent %d bytes to %s\n", n, serverAddr)
}
