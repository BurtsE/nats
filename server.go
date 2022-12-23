package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func newServer() *http.Server {
	http.HandleFunc("/", handler)
	s := http.Server{
		Addr:         "localhost:8008",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	return &s
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("file request")

	id := r.URL.Query()["id"][0]
	m, err := memory.Get(id)
	if err != nil {
		log.Println(err, id)
		w.Write([]byte("message not found"))
	} else {
		rval, _ := json.Marshal(m)
		w.Write(rval)
	}
	log.Println("file sent")
}
