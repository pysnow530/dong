package main

import (
	"flag"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"time"
)

var DB *gorm.DB

var flagAddr = flag.String("addr", "127.0.0.1:5300", "serving address")
var FlagSalt = flag.String("salt", "saltish", "salt for encrypt")
var flagDsn = flag.String("dsn", "root:hello@(localhost)/dong?charset=utf8&parseTime=True", "mysql dsn")

func main() {
	flag.Parse()

	db, err := ConnectDB(*flagDsn)
	if err != nil {
		log.Fatalf("Connect to db failed: %s\n", err.Error())
	}
	DB = db
	defer DB.Close()

	srv := &http.Server{
		Handler:      GetGlobalHandler(),
		Addr:         *flagAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Serving on %s...\n", *flagAddr)
	log.Fatalln(srv.ListenAndServe())
}
