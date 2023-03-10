package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
)

type Repository interface {
	PersistEmails(emails []Mail)
	Search(q string, page int) (*SearchResponse, error)
	GetAll(page int) (*SearchResponse, error)
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
	data, err := toJson(documents)
	if err != nil {
		log.Println(err)
	}
	_, success := r.httpClient.Post("/api/_bulkv2", data)
	if !success {
		log.Println("No se pudo crear la siguiente cantidad de emails: ", len(emails))
	}
}

func (r *ZincRepository) Search(q string, page int) (*SearchResponse, error) {
	query := SearchQuery{Q: q, Page: page}
	path := fmt.Sprintf("/es/%s/_search", r.Index)
	response, success := r.httpClient.Post(path, query.ToJson())
	if !success {
		return nil, logError(fmt.Sprintf("We get an invalid response the query was: %s", string(query.ToJson())))
	}

	searchResponse, err := InitSearchResponse(response, page)
	searchResponse.TotalPages = int(math.Ceil(float64(searchResponse.Hits.Total.Value) / float64(searchResponse.Size)))
	if err != nil {
		return nil, logError(fmt.Sprintf("Error parsing response: %s", string(response)))
	}

	return searchResponse, nil
}

func (r *ZincRepository) GetAll(page int) (*SearchResponse, error) {
	perPage := 30
	data := getAllQuery(page, perPage)
	jsonData, _ := json.Marshal(data)
	path := fmt.Sprintf("/es/%s/_search", r.Index)
	response, success := r.httpClient.Post(path, jsonData)

	if !success {
		return nil, logError(fmt.Sprintf("We get an invalid response the query was: %s", string(jsonData)))
	}

	searchResponse, err := InitSearchResponse(response, page)
	searchResponse.TotalPages = int(math.Ceil(float64(searchResponse.Hits.Total.Value) / float64(searchResponse.Size)))
	if err != nil {
		return nil, logError(fmt.Sprintf("Error parsing response: %s", string(response)))
	}

	return searchResponse, nil
}

func toJson(object interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func logError(message string) error {
	return fmt.Errorf("Algo salio mal")
}

func getAllQuery(page int, perPage int) map[string]interface{} {
	return map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"sort": [1]map[string]string{{"createdAt": "desc"}},
		"size": perPage,
		"from": ((page - 1) * perPage) + 1,
	}
}
