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
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

var memory *Cache

// NOTE: Can test with demo servers.
// nats-sub -s demo.nats.io <subject>

func usage() {
	log.Printf("Usage: nats-sub [-s server] [-creds file] [-nkey file] [-tlscert file] [-tlskey file] [-tlscacert file] [-t] <subject>\n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var userCreds = flag.String("creds", "", "User Credentials File")
	var nkeyFile = flag.String("nkey", "", "NKey Seed File")
	var tlsClientCert = flag.String("tlscert", "", "TLS client certificate file")
	var tlsClientKey = flag.String("tlskey", "", "Private key file for client certificate")
	var tlsCACert = flag.String("tlscacert", "", "CA certificate to verify peer against")
	var showTime = flag.Bool("t", false, "Display timestamps")
	var showHelp = flag.Bool("h", false, "Show help message")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	args := flag.Args()
	if len(args) != 1 {
		showUsageAndExit(1)
	}

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Subscriber")}
	opts = setupConnOptions(opts)

	if *userCreds != "" && *nkeyFile != "" {
		log.Fatal("specify -seed or -creds")
	}

	// Use UserCredentials
	if *userCreds != "" {
		opts = append(opts, nats.UserCredentials(*userCreds))
	}

	// Use TLS client authentication
	if *tlsClientCert != "" && *tlsClientKey != "" {
		opts = append(opts, nats.ClientCert(*tlsClientCert, *tlsClientKey))
	}

	// Use specific CA certificate
	if *tlsCACert != "" {
		opts = append(opts, nats.RootCAs(*tlsCACert))
	}

	// Use Nkey authentication.
	if *nkeyFile != "" {
		opt, err := nats.NkeyOptionFromSeed(*nkeyFile)
		if err != nil {
			log.Fatal(err)
		}
		opts = append(opts, opt)
	}

	memory = NewCashe()

	db := ConnectTODB()
	defer db.Close()
	messages, err := recoverFromDB(db)
	if err != nil {
		log.Println("error while recovering from database: ", err)
	}
	for _, m := range messages {
		memory.Set(m.Order_uid, m)
	}
	// Connect to NATS
	nc, err := nats.Connect(*urls, opts...)
	if err != nil {
		log.Fatal(err)
	}
	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatal("failed stream creation", err)
	}

	subj, i := args[0], 0
	_, err = js.Subscribe(subj, func(msg *nats.Msg) {
		log.Printf("message received on %s\n", subj)
		i += 1
		var m message
		err := json.Unmarshal(msg.Data, &m)

		if err != nil {
			fmt.Println(err)
		} else {
			err = memory.Set(m.Order_uid, m)
			if err != nil {
				log.Println(err)
			}
			addToDB(m, db)

		}
	})
	if err != nil {
		log.Fatal(err)
	}
	s := newServer()
	log.Println("server created")
	s.ListenAndServe()
	log.Println("server launched")

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]", subj)
	if *showTime {
		log.SetFlags(log.LstdFlags)
	}

	//runtime.Goexit()
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}
