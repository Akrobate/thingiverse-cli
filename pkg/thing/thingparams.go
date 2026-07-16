package thing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

const apiBaseURL = "https://api.thingiverse.com"

type ThingResponse struct {
	ID int `json:"id"`
}

type ThingParams struct {
	Id           string `json:"-" yaml:"id"`
	Name         string `json:"name" yaml:"name"`
	Category     string `json:"category" yaml:"category"`
	License      string `json:"license" yaml:"license"`
	IsWip        bool   `json:"is_wip" yaml:"is_wip"`
	Tags         string `json:"tags" yaml:"tags"`
	Instructions string `json:"instructions" yaml:"instructions"`
	Description  string `json:"description" yaml:"description"`
}

func NewThingParams() (*ThingParams, error) {

	return &ThingParams{}, nil
}

func (tp *ThingParams) Save() error {
	data, err := yaml.Marshal(tp)
	if err != nil {
		return err
	}

	return os.WriteFile("./thingiverse.yml", data, 0644)
}

func (tp *ThingParams) Create(accessToken string) (int, error) {

	jsonData, err := json.Marshal(tp)
	if err != nil {
		return 0, fmt.Errorf("Error JSON serialize : %w", err)
	}

	url := fmt.Sprintf("%s/things", apiBaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("Error creating request : %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("Request failed : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("API Error (HTTP %d) : %s", resp.StatusCode, string(bodyBytes))
	}

	var thingResp ThingResponse
	if err := json.NewDecoder(resp.Body).Decode(&thingResp); err != nil {
		return 0, fmt.Errorf("Parse response problem : %w", err)
	}

	tp.Id = fmt.Sprintf("%d", thingResp.ID)

	return thingResp.ID, nil
}
