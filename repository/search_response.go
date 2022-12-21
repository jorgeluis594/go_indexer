package repository

import "encoding/json"

type SearchResponse struct {
	CurrentPage int
	Size        int
	Hits        Hits `json:"hits"`
}

type Hits struct {
	Total HitsTotal `json:"total"`
	Hits  []Hit     `json:"hits"`
}

type Hit struct {
	Source Mail `json:"_source"`
}

type HitsTotal struct {
	Value int `json:"value"`
}

func InitSearchResponse(response []byte, currentPage int, size int) (*SearchResponse, error) {
	searchResponse := SearchResponse{CurrentPage: currentPage, Size: size}
	err := json.Unmarshal(response, &searchResponse)
	if err != nil {
		return nil, err
	}
	return &searchResponse, nil
}
