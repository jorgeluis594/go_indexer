package indexer

import (
	"fmt"
	"io"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

type Mail struct {
	MailId          string            `json:"mailId"`
	Date            time.Time         `json:"date"`
	EmailSender     string            `json:"emailSender"`
	EmailReceivers  []string          `json:"emailReceivers"`
	CopiedReceivers []string          `json:"copiedReceivers"`
	HiddenReceivers []string          `json:"hiddenReceivers"`
	Subject         string            `json:"subject"`
	CustomHeaders   map[string]string `json:"customHeaders"`
	Content         string            `json:"content"`
}

func InitMail(mailReader *mail.Message) (*Mail, error) {
	headers := mailReader.Header
	newMail := Mail{
		MailId:          cleanScapeCharacters(headers.Get("Message-Id")),
		EmailSender:     headers.Get("From"),
		EmailReceivers:  strings.Split(headers.Get("To"), ", "),
		CopiedReceivers: strings.Split(headers.Get("Cc"), ", "),
		HiddenReceivers: strings.Split(headers.Get("Bcc"), ", "),
		Subject:         headers.Get("Subject"),
	}

	var err error
	newMail.Date, err = parseDate(headers.Get("Date"))
	newMail.CustomHeaders, err = customHeadersFor(headers)
	newMail.Content, err = readBody(mailReader.Body)

	if err != nil {
		return &newMail, err
	}

	return &newMail, nil
}

func cleanScapeCharacters(info string) string {
	id := strings.Replace(info, "<", "", 1)
	return strings.Replace(id, ">", "", 1)
}

func parseDate(date string) (time.Time, error) {
	return time.Parse(time.RFC822, date)
}

func customHeadersFor(headers mail.Header) (map[string]string, error) {
	var customHeaders = make(map[string]string)
	for key, value := range headers {
		matched, err := regexp.MatchString("^X-", key)
		if err != nil {
			return make(map[string]string), fmt.Errorf("error comparing key from headers: %s", key)
		}

		if matched {
			customHeaders[key] = value[0]
		}
	}

	return customHeaders, nil
}

func readBody(body io.Reader) (string, error) {
	content, err := io.ReadAll(body)
	if err != nil {
		return "", fmt.Errorf("cannot read body for email")
	}
	return string(content), nil
}
