// Copyright 2012-2022 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"main/server"
	"main/storage"
)

func main() {

	server.Createsub("nats-sub")
	storage.ConnectTODB()
	db := storage.GetDatabase()
	defer db.Close()
	err := storage.RecoverFromDB()
	log.Println("created cache")
	if err != nil {
		log.Println("error while recovering from database: ", err)
	}
	s := server.NewServer()
	log.Println("server created")
	err = s.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
