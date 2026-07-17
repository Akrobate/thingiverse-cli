package thing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ThingUpdateRequest struct {
	Name         string   `json:"name" yaml:"name"`
	Category     string   `json:"category" yaml:"category"`
	License      string   `json:"license" yaml:"license"`
	IsWip        int      `json:"is_wip" yaml:"is_wip"`
	Tags         []string `json:"tags" yaml:"tags"`
	Instructions string   `json:"instructions" yaml:"instructions"`
	Description  string   `json:"description" yaml:"description"`
}

func Update(id int, tr *ThingUpdateRequest, accessToken string) error {

	jsonData, err := json.Marshal(tr)
	if err != nil {
		return fmt.Errorf("Error JSON serialize : %w", err)
	}

	url := fmt.Sprintf("%s/things/%d", apiBaseURL, id)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error creating request : %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Request failed : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API Error (HTTP %d) : %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
