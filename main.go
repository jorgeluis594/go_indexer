package main

import (
	"fmt"
	"log"
	"net/mail"
	"os"
)

func main() {

	file, err := os.Open("1.")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	msg, err := mail.ReadMessage(file)
	if err != nil {
		log.Fatal(err)
	}

	// Imprime algunos datos del mensaje
	fmt.Println(msg.Header)
}
