package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
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
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server")
		return
	}
	defer ln.Close()

	fmt.Println("Server started on port 8080")

	sema := make(chan struct{}, 200)

	connId := 0

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}
		fmt.Println("Connection accepted", connId)
		go handleConnection(conn, sema, connId)
		connId++
	}

}

func handleConnection(conn net.Conn, sema chan struct{}, connId int) {
	defer conn.Close()
	defer fmt.Println("Connection closed", connId)

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	inputCh := make(chan Record)
	outputCh := make(chan ProcessedRecord)
	workerWg := &sync.WaitGroup{}

	stopCh := make(chan struct{})

	go func() {
		workers := 100
		for i := 0; i < workers; i++ {
			sema <- struct{}{}
			workerWg.Add(1)
			go worker(inputCh, outputCh, sema, workerWg, stopCh)
		}
	}()

	replyWg := sync.WaitGroup{}
	replyWg.Add(1)

	go func() {
		defer replyWg.Done()
		for processedRecord := range outputCh {
			err := encoder.Encode(processedRecord)
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				close(stopCh)
				return
			}
		}
	}()

	for {
		var records []Record
		err := decoder.Decode(&records)
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println("Error decoding JSON:", err)
			close(stopCh)
			return
		}

		for _, record := range records {
			inputCh <- record
		}
	}

	close(inputCh)
	workerWg.Wait()

	close(outputCh)
	replyWg.Wait()
}

func worker(inputCh chan Record, outputCh chan ProcessedRecord, sema chan struct{}, wg *sync.WaitGroup, stopCh chan struct{}) {
	defer func() {
		wg.Done()
		<-sema
	}()

	for {
		select {
		case record, ok := <-inputCh:
			if !ok {
				return
			}
			select {
			case outputCh <- doHttpRequest(record):
			case <-stopCh:
				return
			}
		case <-stopCh:
			return
		}
	}

}

func doHttpRequest(record Record) ProcessedRecord {
	client := &http.Client{}
	reqBody, _ := json.Marshal(record)
	req, _ := http.NewRequest("POST", "http://localhost:3000", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return ProcessedRecord{Record: record, Status: "Failed", Response: err.Error()}
	}
	defer resp.Body.Close()
	var processedRecord ProcessedRecord = ProcessedRecord{Record: record, Status: "Success"}
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		processedRecord.Status = "Failed"
		processedRecord.Response = err.Error()
	}
	processedRecord.Response = string(resBody)

	return processedRecord
}
