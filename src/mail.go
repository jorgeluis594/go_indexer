package indexer

import "time"

type Mail struct {
	Id               string            `json:"_id"`
	Date             time.Time         `json:"date"`
	EmailSender      string            `json:"emailSender"`
	EmailReceivers   []string          `json:"emailReceivers"`
	CopiedReceivers  []string          `json:"copiedReceivers"`
	HiddenReceivers  []string          `json:"hiddenReceivers"`
	Subject          string            `json:"subject"`
	CustomAttributes map[string]string `json:"customAttributes"`
	Content          string            `json:"content"`
}
