package thing

import (
	"encoding/json"
	"fmt"

	"github.com/Akrobate/thingiverse-cli/pkg/utils"
)

type CategoryGetResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
	Slug  string `json:"slug"`
}

type SubCategoryGetResponse struct {
	Children []CategoryGetResponse `json:"children"`
}

func CategorySearch(accessToken string) (*[]CategoryGetResponse, error) {

	url := fmt.Sprintf("%s/categories", apiBaseURL)
	resp, err := utils.HttpDoAuthenticatedGetRequest(url, accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var t []CategoryGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("parse response problem: %w", err)
	}

	return &t, nil
}

func SubCategorySearch(slug string, accessToken string) (*SubCategoryGetResponse, error) {

	url := fmt.Sprintf("%s/categories/%s", apiBaseURL, slug)
	resp, err := utils.HttpDoAuthenticatedGetRequest(url, accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var t SubCategoryGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("parse response problem: %w", err)
	}

	return &t, nil
}
