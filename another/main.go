package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Record struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProcessedRecord struct {
	Record   Record `json:"record"`
	Status   string `json:"status"`
	Response string `json:"response"`
}

func main() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("Error starting server")
		return
	}
	defer ln.Close()

	fmt.Println("Server started on port 8081")

	connId := 0

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}
		fmt.Println("Connection accepted", connId)
		go handleConnection(conn, connId)
		connId++
	}

}

func handleConnection(conn net.Conn, connId int) {
	defer conn.Close()
	defer fmt.Println("Connection closed", connId)

	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()
	logger := log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	decoder := json.NewDecoder(conn)

	for {
		fmt.Println("Reading from connection")
		logger.Println("Reading from connection")
		var records []Record
		err := decoder.Decode(&records)
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println("Error decoding JSON:", err)
			return
		}

		for _, record := range records {
			logger.Println("Record received", record)
		}
	}

}
