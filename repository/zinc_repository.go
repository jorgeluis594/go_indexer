package repository

import (
	"encoding/json"
	"fmt"
	"log"
)

type Repository interface {
	PersistEmails(emails []Mail)
	Search(q string, page int) (*SearchResponse, error)
}

type ZincRepository struct {
	httpClient Http
	Index      string
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
	documents := documentsBulk{Index: r.Index, Records: emails}
	_, success := r.httpClient.Post("/api/_bulkv2", toJson(documents))
	if !success {
		log.Println("No se pudo crear la siguiente cantidad de emails: ", len(emails))
	}
}

func (r *ZincRepository) Search(q string, page int) (*SearchResponse, error) {
	query := SearchQuery{Q: q, Page: page}
	path := fmt.Sprintf("/es/%s/_search", r.Index)
	response, success := r.httpClient.Post(path, query.ToJson())
	if !success {
		var err Error
		json.Unmarshal(response, &err)
		return nil, fmt.Errorf("request error: %s", err.Message)
	}

	var searchResponse SearchResponse
	err := json.Unmarshal(response, &searchResponse)
	if err != nil {
		return nil, err
	}

	return &searchResponse, nil
}

func toJson(object interface{}) []byte {
	jsonData, err := json.Marshal(object)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData
}
