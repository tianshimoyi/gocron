package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Request struct {
	Msg Message `json:"message"`
}

type Message struct {
	Sid     string   `json:"sender_id"`
	Rids    []string `json:"receiver_ids"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Cid     string   `json:"title"`
}

func main() {
	http.HandleFunc("/message", MockNotify)
	http.ListenAndServe("0.0.0.0:9090", nil)
}

func MockNotify(writer http.ResponseWriter, request *http.Request) {
	log.Println("receive a request")
	var req Request
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		request.Body.Close()
		log.Fatal(err)
	}
	log.Println(req)
	writer.WriteHeader(http.StatusAccepted)
}
