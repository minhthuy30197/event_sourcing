package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func main() {
	var requests []*http.Request
	var clients []*http.Client
	for i := 0; i < 10; i++ {
		var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
		request, err := http.NewRequest("POST", "http://localhost:8080/es/test-version", bytes.NewBuffer(jsonStr))
		request.Header.Set("X-Custom-Header", "myvalue")
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			panic(err)
		}
		requests = append(requests, request)
		client := &http.Client{}
		clients = append(clients, client)
	}
	log.Println("- ", requests)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go SendReq(clients[i], requests[i], i)
	}
	wg.Wait()
}

func SendReq(client *http.Client, request *http.Request, i int) {
	resp, err := client.Do(request)
	if err != nil {
		log.Println("---- : ", i, " ", err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	defer resp.Body.Close()
}
