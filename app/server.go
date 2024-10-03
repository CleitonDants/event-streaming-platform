package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

// 	Uncomment this block to pass the first stage

	listener, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}else{
		fmt.Println("Server started on port 9092")
	}

	conn, err := listener.Accept()
	
	// if conn
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	} else {
		fmt.Println("Accepted connection: ", (conn.LocalAddr().String()))
	}

	defer conn.Close()

	request := make([]byte, 1024)
	conn.Read(request)

	request_length := int32(binary.BigEndian.Uint32(request[0:4]))
	request_api_key := int16(binary.BigEndian.Uint16(request[4:6]))
	request_api_version := int16(binary.BigEndian.Uint16(request[6:8]))
	correlation_id := int32(binary.BigEndian.Uint32(request[8:12]))
	client_id := string(request[12:])

	fmt.Printf("request length: %d\n", request_length) //4 bytes
	fmt.Printf("request_api_key: %d\n", request_api_key) //2 bytes
	fmt.Printf("request_api_version: %d\n", request_api_version) //2 bytes
	fmt.Printf("correlation_id: %d\n", correlation_id) //4 bytes
	fmt.Printf("client_id: %s\n", client_id)

	response := make([]byte, 8)
	copy(response[4:], request[8:12])

	conn.Write(response)
}
