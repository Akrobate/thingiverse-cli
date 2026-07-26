package thing

import (
	"encoding/json"
	"fmt"

	"github.com/Akrobate/thingiverse-cli/pkg/utils"
)

type ThingGetResponse struct {
	Id           int                   `json:"id"`
	Name         string                `json:"name"`
	Category     string                `json:"category"`
	License      string                `json:"license"`
	IsWip        int                   `json:"is_wip"`
	Tags         []ThingTagGetResponse `json:"tags"`
	Instructions string                `json:"instructions"`
	Description  string                `json:"description"`
}

type ThingTagGetResponse struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

func Get(id string, accessToken string) (*ThingGetResponse, error) {

	url := fmt.Sprintf("%s/things/%s", apiBaseURL, id)

	resp, err := utils.HttpDoAuthenticatedGetRequest(url, accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var t ThingGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("parse response problem: %w", err)
	}

	return &t, nil
}
