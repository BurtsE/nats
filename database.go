package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type config struct {
	Host     string `json:host`
	Port     int    `json:port`
	User     string `json:user`
	Password string `json:password`
	Dbname   string `json:dbname`
}

func ConnectToDB() {
	log.Println("connecting to db...")
	var conf config
	file, err := os.Open("config.json")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	cdata, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(cdata, &conf)
	if err != nil {
		log.Fatal(err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Dbname)

	db, err := sql.Open("pgx", psqlInfo)

	if err != nil {
		log.Printf("unable to connect to db", err)
		os.Exit(1)
	}
	defer db.Close()
	var data *sql.Rows
	data, err = db.Query("select * from orders")
	if err != nil {
		log.Println(err)
	}
	for data.Next() {
		var s []byte
		data.Scan(&s)
		log.Println(string(s))
	}
	log.Println("connected to db")
}
