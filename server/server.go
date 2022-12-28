package server

import (
	"log"
	"main/storage"
	"net/http"
	"text/template"
	"time"
)

func NewServer() *http.Server {
	http.HandleFunc("/", handler)
	s := http.Server{
		Addr:         "localhost:8000",
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
	t, err := template.New("index.html").ParseFiles("../static/index.html")
	if err != nil {
		log.Println(err)
	}
	memory := storage.GetCache()

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
