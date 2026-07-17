package thing

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API Error (HTTP %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var t ThingGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("parse response problem: %w", err)
	}

	return &t, nil
}
