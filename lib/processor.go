package lib

import (
	"bytes"
	"github.com/jorgeluis594/go_indexer/repository"
	"log"
	"math"
	"net/mail"
	"os"
)

type Processor struct {
	paths            []string
	emailDataChannel chan []byte
	emailsChannel    chan *repository.Mail
	repository       repository.Repository
	emailsToSend     []repository.Mail
	emailsProcessed  int
}

func InitProcessor(filePaths []string, r repository.Repository) *Processor {
	return &Processor{
		paths:            filePaths,
		emailDataChannel: make(chan []byte),
		emailsChannel:    make(chan *repository.Mail),
		repository:       r,
	}
}

func (p *Processor) Process() {
	// reading and parsing emails
	for _, path := range p.paths {
		go p.readFile(path)
	}

	for i := 0; i < len(p.paths)*2; i++ {
		select {
		case emailData := <-p.emailDataChannel:
			go p.parseEmail(emailData)
		case email := <-p.emailsChannel:
			if email == nil {
				continue
			}
			p.emailsToSend = append(p.emailsToSend, *email)
		}
	}
	close(p.emailsChannel)
	close(p.emailDataChannel)

	numOfRequest := int(math.Ceil(float64(len(p.emailsToSend)) / 1000))
	log.Println(numOfRequest)
	for i := 0; i < numOfRequest; i++ {
		var availableEmailsToSend int
		if len(p.emailsToSend) < 1000 {
			availableEmailsToSend = len(p.emailsToSend)
		} else {
			availableEmailsToSend = 1000
		}
		emailsToSend := p.emailsToSend[:availableEmailsToSend]
		p.emailsToSend = p.emailsToSend[1000:]
		log.Println("cantidad de emails a enviar", len(emailsToSend))
		log.Println("cantidad de emails a restantes", len(p.emailsToSend))
		p.persistEmails(emailsToSend)
	}
}

func (p *Processor) readFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println("Error reading file: ", path)
		p.emailDataChannel <- nil
	}

	p.emailDataChannel <- data
}

func (p *Processor) parseEmail(data []byte) {
	if data == nil {
		p.emailsChannel <- nil
		return
	}

	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		p.emailsChannel <- nil
		return
	}

	email, err := repository.InitMail(msg)
	if err != nil {
		p.emailsChannel <- nil
		return
	}

	p.emailsChannel <- email
}

func (p *Processor) persistEmails(emails []repository.Mail) {
	p.repository.PersistEmails(emails)
}
