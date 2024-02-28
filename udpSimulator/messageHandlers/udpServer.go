package messageHandlers

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
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
		pushResponse("configResponse", "Configuration successful!")
		_, err = conn.WriteToUDP(response, addr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		}
	}
}

// pushResponse
func pushResponse(topic string, payload string) {
	//kafkaService := "kafka-service.default.svc.cluster.local:9099"
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://kafka-service.default.svc.cluster.local:9099/kafkasvc/msg/send"), strings.NewReader(payload))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
	}

	req.Header.Set("message-topic", topic)
	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making HTTP request:", err)
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
	}
	log.Println("Response pushed to kafka successfully")
}
