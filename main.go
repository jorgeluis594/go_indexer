package main

import (
	"bytes"
	"fmt"
	indexer "github.com/jorgeluis594/go_indexer/src"
	"log"
	"net/mail"
	"os"
)

func main() {
	path := "harris-s"
	emails, success := loadEmails(path)
	if success {
		fmt.Println(len(*emails))
	}
}

func loadEmails(path string) (*[]indexer.Mail, bool) {
	emails := make([]indexer.Mail, 0)
	directory, err := indexer.InitDirectory("harris-s")

	if err != nil {
		log.Fatal("Error reading directory: ", path)
		return &emails, false
	}

	for _, path := range directory.GetPaths() {
		emailReader, success := readEmail(path)
		if !success {
			continue
		}
		email, err := indexer.InitMail(emailReader)
		if err != nil {
			log.Printf("Error parsing email with path: %s\n", path)
			continue
		}
		emails = append(emails, *email)
	}

	return &emails, true
}

func readEmail(path string) (*mail.Message, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println("Error reading file: ", path)
		return nil, false
	}

	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		log.Printf("File with path: %s is not a email\n", path)
		return nil, false
	}

	return msg, true
}
