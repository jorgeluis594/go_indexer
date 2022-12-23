package repository

import (
	"encoding/json"
	"log"
)

type SearchQuery struct {
	Q    string
	Page int
}

func (q *SearchQuery) ToJson() []byte {
	perPage := 30
	data := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{"term": map[string]string{"copiedReceivers": q.Q}},
					{"match_phrase": map[string]string{"content": q.Q}},
				},
			},
		},
		"sort": [1]map[string]string{{"createdAt": "asc"}},
		"size": perPage,
		"from": ((q.Page - 1) * perPage) + 1,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Cannot marshal next data: %v", data)
	}

	return jsonData
}
