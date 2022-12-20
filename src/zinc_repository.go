package indexer

import "log"

type ZincRepository struct {
	httpClient Http
}

type Error struct {
	Message     string `json:"mailId"`
	RecordCount int    `json:"record_count"`
}

type documentsBulk struct {
	Index   string `json:"index"`
	Records []Mail `json:"records"`
}

func InitRepository(httpClient Http) *ZincRepository {
	return &ZincRepository{httpClient: httpClient}
}

func (r *ZincRepository) PersistEmails(emails []Mail) {
	documents := documentsBulk{Index: "emails", Records: emails}
	_, success := r.httpClient.Post("/api/_bulkv2", documents)
	if !success {
		log.Println("No se pudo crear la siguiente cantidad de emails: ", len(emails))
	}
}
