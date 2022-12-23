package lib

import (
	"bytes"
	"fmt"
	"github.com/jorgeluis594/go_indexer/repository"
	"log"
	"net/mail"
	"os"
)

type Processor struct {
	paths         []string
	pathsChannel  chan string
	emailsChannel chan *repository.Mail
	doneChannel   chan bool
	repository    repository.Repository
	emailsToSend  []repository.Mail
}

func (p *Processor) InitProcessor(filePaths []string) {
	p.paths = filePaths
	p.pathsChannel = make(chan string)
	p.emailsChannel = make(chan *repository.Mail)
	p.doneChannel = make(chan bool)
}

func (p *Processor) Process() {
	p.sendPathsToChannel()
	for _, path := range p.paths {
		go p.readFile(path, p.processEmail)
	}
	p.ensureToCloseRoutines()
	close(p.emailsChannel)
	close(p.pathsChannel)
	close(p.doneChannel)
}

func (p *Processor) sendPathsToChannel() {
	for _, path := range p.paths {
		p.pathsChannel <- path
	}
}

func (p *Processor) ensureToCloseRoutines() {
	for i := 0; i < len(p.paths); i++ {
		<-p.doneChannel
	}
}

func (p *Processor) readFile(path string, callback func(string, []byte, func(*repository.Mail))) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println("Error reading file: ", path)
		p.doneChannel <- true
	}
	go callback(path, data, p.persistEmails)
}

func (p *Processor) processEmail(path string, data []byte, callback func(*repository.Mail)) {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		log.Printf("File with path: %s is not a email\n", path)
		p.doneChannel <- true
	}

	email, err := repository.InitMail(msg)
	if err != nil {
		log.Printf("Error parsing email with path: %s\n", path)
		p.doneChannel <- true
	}

	go callback(email)
}

func (p *Processor) persistEmails(email *repository.Mail) {
	p.emailsToSend = append(p.emailsToSend, *email)
	if len(p.emailsToSend) == 1000 || len(p.doneChannel)+1 == len(p.paths) {
		fmt.Printf("Sending %d emails", len(p.emailsToSend))
		p.repository.PersistEmails(p.emailsToSend)
	}

	p.doneChannel <- true
}
