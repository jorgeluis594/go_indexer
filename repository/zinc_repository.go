package repository

import (
	"encoding/json"
	"fmt"
	"log"
)

type Repository interface {
	PersistEmails(emails []Mail)
	Search(q string, page int) *SearchResponse
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

func InitRepository(httpClient Http, index string) *ZincRepository {
	return &ZincRepository{httpClient: httpClient, Index: index}
}

func (r *ZincRepository) PersistEmails(emails []Mail) {
	documents := documentsBulk{Index: r.Index, Records: emails}
	_, success := r.httpClient.Post("/api/_bulkv2", toJson(documents))
	if !success {
		log.Println("No se pudo crear la siguiente cantidad de emails: ", len(emails))
	}
}

func (r *ZincRepository) Search(q string, page int) *SearchResponse {
	query := SearchQuery{Q: q, Page: page}
	path := fmt.Sprintf("/es/%s/_search", r.Index)
	response, success := r.httpClient.Post(path, query.ToJson())
	if !success {
		log.Fatalf("We get an invalid response the query was: %v", query)
	}

	searchResponse, err := InitSearchResponse(response, page)
	if err != nil {
		log.Fatalf("Error parsing response: %s", string(response))
	}

	return searchResponse
}

func toJson(object interface{}) []byte {
	jsonData, err := json.Marshal(object)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData
}
