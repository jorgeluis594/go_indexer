package main

import (
	"encoding/json"
	"fmt"
	indexer "github.com/jorgeluis594/go_indexer/src"
	"log"
	"net/mail"
	"os"
)

func main() {

	file, err := os.Open("4.")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	msg, err := mail.ReadMessage(file)
	if err != nil {
		log.Fatal(err)
	}

	email, err := indexer.InitMail(msg)
	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.Marshal(email)
	if err != nil {
		log.Fatal("Error parsing email to json")
	}

	fmt.Println(string(jsonData))
}
