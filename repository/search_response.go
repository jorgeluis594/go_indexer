package repository

import (
	"encoding/json"
)

type SearchResponse struct {
	CurrentPage int  `json:"currentPage"`
	Size        int  `json:"size"`
	TotalPages  int  `json:"totalPages"`
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

func InitSearchResponse(response []byte, currentPage int) (*SearchResponse, error) {
	searchResponse := SearchResponse{CurrentPage: currentPage, Size: 30}
	err := json.Unmarshal(response, &searchResponse)
	if err != nil {
		return nil, err
	}
	return &searchResponse, nil
}
