package main

import (
	"bytes"
	"flag"
	"fmt"
	indexer "github.com/jorgeluis594/go_indexer/src"
	"log"
	"net/mail"
	"os"
)

func main() {
	path := flag.String("path", "", "path to index")
	host := flag.String("host", "", "host of Zinc Search client")
	username := flag.String("username", "", "username of db")
	password := flag.String("password", "", "password of db")
	flag.Parse()
	
	clientHttp := indexer.InitHttpClient(*host, *username, *password)
	repository := indexer.InitRepository(clientHttp)
	emails, success := loadEmails(*path)
	if !success {
		log.Fatal("No se pudieron cargar los emails del directorio: ", path)
	}

	for len(emails) > 0 {
		var numberOfEmails int
		if len(emails) < 1000 {
			numberOfEmails = len(emails)
		} else {
			numberOfEmails = 1000
		}
		emailsToSend := emails[:numberOfEmails]
		emails = emails[numberOfEmails:]
		fmt.Println("Count of sent emails: ", len(emailsToSend))
		repository.PersistEmails(emailsToSend)
	}
}

func loadEmails(path string) ([]indexer.Mail, bool) {
	emails := make([]indexer.Mail, 0)
	directory, err := indexer.InitDirectory("harris-s")

	if err != nil {
		log.Fatal("Error reading directory: ", path)
		return emails, false
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

	return emails, true
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
