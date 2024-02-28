package messageHandlers

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

// Message struct definition
type Message struct {
	MessageType     int    `json:"MessageType"`
	ActualCommand   int    `json:"ActualCommand"`
	SenderID        uint32 `json:"SenderID"`
	ReceiverID      uint32 `json:"ReceiverID"`
	MsgID           uint32 `json:"MsgID"`
	SequenceNumber  uint32 `json:"SequenceNumber"`
	TimestampSender int64  `json:"TimestampSender"`
	Payload         string `json:"Payload"`
}

// SendMessage
func SendMessage(ctx *gin.Context) {
	tr := otel.Tracer("udp-service")
	ctx1 := context.Background()
	ctx1, span := tr.Start(ctx1, "udp-svc")
	defer span.End()

	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	defer conn.Close()
	fmt.Println("Sending message")

	reqBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading request body"})
		return
	}

	var message Message
	err = json.Unmarshal(reqBody, &message)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	buffer := new(bytes.Buffer)
	payloadBytes := []byte(message.Payload)

	binary.Write(buffer, binary.LittleEndian, uint32(len(payloadBytes)))

	buffer.WriteByte(byte(message.MessageType))   // msgType 1 - CMD
	buffer.WriteByte(byte(message.ActualCommand)) // msgCommand 1- START

	binary.Write(buffer, binary.LittleEndian, uint32(message.SenderID))       // Sender ID
	binary.Write(buffer, binary.LittleEndian, uint32(message.ReceiverID))     // Receiver ID
	binary.Write(buffer, binary.LittleEndian, uint32(message.MsgID))          // Message ID
	binary.Write(buffer, binary.LittleEndian, uint32(message.SequenceNumber)) // Sequence Number
	binary.Write(buffer, binary.LittleEndian, message.TimestampSender)        // Timestamp Sender

	// Set payload itself
	buffer.Write(payloadBytes)

	n, err := conn.Write(buffer.Bytes())
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}
	fmt.Printf("Sent %d bytes to %s\n", n, serverAddr)

	responseBuffer := make([]byte, 1024) // Adjust the buffer size based on your requirements
	n, _, err = conn.ReadFromUDP(responseBuffer)
	if err != nil {
		fmt.Println("Error receiving data:", err)
		return
	}

	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}
	ctx.JSON(200, string(responseBuffer[:n]))
}
