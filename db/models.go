package db

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ModelRecord struct {
	CollectionID   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	Created        string `json:"created"`
	Updated        string `json:"updated"`
}

type ModelRecordsResponse struct {
	Page       int           `json:"page"`
	PerPage    int           `json:"perPage"`
	TotalPages int           `json:"totalPages"`
	TotalItems int           `json:"totalItems"`
	Items      []ModelRecord `json:"items"`
}


const (
	QueryParamPage      QueryParam = "page"
	QueryParamPerPage   QueryParam = "perPage"
	QueryParamSort      QueryParam = "sort"
	QueryParamFilter    QueryParam = "filter"
	QueryParamSkipTotal QueryParam = "skipTotal"
)

func (db *db) GetModelRecords(queryParams map[QueryParam]string) (*ModelRecordsResponse, error) {
	recordsEndpoint := "/api/collections/models/records"

	queryString := ""
	if len(queryParams) > 0 {
		query := url.Values{}
		for key, value := range queryParams {
			query.Add(string(key), value)
		}
		queryString = "?" + query.Encode()
	}

	req, err := http.NewRequest("GET", db.Url+recordsEndpoint+queryString, nil)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := db.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[ERROR] request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var recordsResponse ModelRecordsResponse
	if err := json.Unmarshal(respBody, &recordsResponse); err != nil {
		return nil, fmt.Errorf("[ERROR] failed to parse response: %w", err)
	}

	return &recordsResponse, nil
}

func GroupModelRecordsByCompany(records []ModelRecord) map[string][]ModelRecord {
	groupedModels := make(map[string][]ModelRecord)

	for _, record := range records {
		parts :=  strings.Split(record.Name, "/")
		if len(parts) != 2 {
			continue
		}
		company := parts[0]
		groupedModels[company] = append(groupedModels[company], record)
	}

	return groupedModels
}
