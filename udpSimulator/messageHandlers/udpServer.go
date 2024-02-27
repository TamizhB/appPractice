package messageHandlers

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// HandleMessage
func HandleMessage(data []byte) {
	buffer := bytes.NewReader(data)

	var payloadSize uint32
	binary.Read(buffer, binary.LittleEndian, &payloadSize)

	msgType, _ := buffer.ReadByte()
	msgCommand, _ := buffer.ReadByte()

	var senderID, receiverID, msgID, sequenceNumber uint32
	var timestampSender int64

	binary.Read(buffer, binary.LittleEndian, &senderID)
	binary.Read(buffer, binary.LittleEndian, &receiverID)
	binary.Read(buffer, binary.LittleEndian, &msgID)
	binary.Read(buffer, binary.LittleEndian, &sequenceNumber)
	binary.Read(buffer, binary.LittleEndian, &timestampSender)

	payloadData := make([]byte, payloadSize)
	_, err := buffer.Read(payloadData)
	if err != nil {
		fmt.Println("Error reading payload data:", err)
		return
	}

	fmt.Printf("Payload Size: %d\n", payloadSize)
	fmt.Printf("Message Type: %d\n", msgType)
	fmt.Printf("Message Command: %d\n", msgCommand)
	fmt.Printf("Sender ID: %d\n", senderID)
	fmt.Printf("Receiver ID: %d\n", receiverID)
	fmt.Printf("Message ID: %d\n", msgID)
	fmt.Printf("Sequence Number: %d\n", sequenceNumber)
	fmt.Printf("Timestamp Sender: %v\n", time.Unix(0, timestampSender))
	fmt.Printf("Payload: %v\n", string(payloadData))
}

// Serve
func Serve() {
	listenAddr, err := net.ResolveUDPAddr("udp", ":8888")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}

	// defer conn.Close() ---> runs indefinitely

	buffer := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading data:", err)
			continue
		}

		fmt.Printf("Server Received %d bytes from %s:", n, addr)
		HandleMessage(buffer[:n])
		response := []byte("Configuration successful!")
		_, err = conn.WriteToUDP(response, addr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		}
	}
}
