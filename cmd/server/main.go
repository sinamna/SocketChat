package main

import (
	"github.com/sinamna/SocketChat/internal/handlers"
	"log"
	"net/http"
)

func main(){
	mux := routes()

	log.Println("starting channel listener")
	go handlers.ListenToPayloadChan()

	log.Println("starting server on port 8080")
	_ = http.ListenAndServe(":8080", mux)
}