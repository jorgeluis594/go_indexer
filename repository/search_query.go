package repository

import "encoding/json"

type SearchQuery struct {
	Q    string
	Page int
}

func (q *SearchQuery) ToJson() ([]byte, error) {
	perPage := 30
	data := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{"term": map[string]string{"copiedReceivers": q.Q}},
				},
			},
		},
		"size": perPage,
		"from": (q.Page - 1) * perPage,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
