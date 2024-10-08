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

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Accepted connection: ", (conn.LocalAddr().String()))
	request := make([]byte, 1024)

	for {
		n, err := conn.Read(request)
		if err != nil {
			fmt.Println("Reading error:", err)
			return
		}
		fmt.Printf("Received: %s\n", string(request[:n]))

		// Enviar resposta ao cliente (eco de volta)
		_, err = conn.Write(kafkaProtocol(request))
		if err != nil {
			fmt.Println("Sending response error:", err)
			return
		}
	}
	
}

func initListener() net.Listener {
	listener, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	
	fmt.Println("Server started on port 9092")
	
	return listener
}

func kafkaProtocol(request []byte) []byte {

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
	
	if request_api_version < 0 || request_api_version > 4{
		fmt.Println("error: Unsupported Kafka API version")
		response = append(response,0,byte(35))
        return response
	}

	response = binary.BigEndian.AppendUint16(response, 0)
	response = append(response, 2)
	response = binary.BigEndian.AppendUint16(response, 18)
	response = binary.BigEndian.AppendUint16(response, 0)
	response = binary.BigEndian.AppendUint16(response, 4)
	response = append(response, 0)
	response = binary.BigEndian.AppendUint32(response, 0)
	response = append(response, 0)
	binary.BigEndian.PutUint32(response[0:4], uint32(len(response)-4))
	
	fmt.Printf("response: %d\n", response)
	return response
}

func main() {
	listener := initListener()
	defer listener.Close()
	
	for {

		conn,error := listener.Accept()
		if error != nil{
			fmt.Println("Error accepting connection: ", error.Error())
            continue
		}

		go handleConnection(conn)
	}

}
