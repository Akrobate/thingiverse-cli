package thing

import (
	"encoding/json"
	"fmt"

	"github.com/Akrobate/thingiverse-cli/pkg/utils"
)

type TagGetResponse struct {
	Tag      string `json:"tag"`
	TagCount int    `json:"tag_count"`
}

func TagSearch(searchString string, accessToken string) (*[]TagGetResponse, error) {

	url := fmt.Sprintf("%s/tags/%s/search-tags", apiBaseURL, searchString)
	resp, err := utils.HttpDoAuthenticatedGetRequest(url, accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var t []TagGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("parse response problem: %w", err)
	}

	return &t, nil
}
