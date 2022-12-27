package main

import (
	"log"
	"net/http"
	"text/template"
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
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	id := r.Form.Get("id")

	t, err := template.New("index.html").ParseFiles("index.html")

	m, err := memory.Get(id)
	if err != nil {
		log.Println(err, id)
		err = t.Execute(w, "message not found")

	} else {

		err = t.Execute(w, m)
	}
	if err != nil {
		log.Println(err)
	}
	log.Println("file sent")
}
