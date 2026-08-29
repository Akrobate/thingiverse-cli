package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func HttpDoAuthenticatedRequest(url string, method string, body io.Reader, accessToken string) (*http.Response, error) {
	return HttpDoAuthenticatedRequestWithContentType(url, method, body, "application/json", accessToken)
}

func HttpDoAuthenticatedRequestWithContentType(url string, method string, body io.Reader, contentType string, accessToken string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

func HttpDoAuthenticatedGetRequest(url string, accessToken string) (*http.Response, error) {
	resp, err := HttpDoAuthenticatedRequest(url, "GET", nil, accessToken)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API GET Error (HTTP %d): %s", resp.StatusCode, string(bodyBytes))
	}
	return resp, err
}

func HttpDoAuthenticatedPostRequest(url string, jsonData []byte, accessToken string) (*http.Response, error) {
	resp, err := HttpDoAuthenticatedRequest(url, "POST", bytes.NewBuffer(jsonData), accessToken)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API POST Error (HTTP %d): %s", resp.StatusCode, string(bodyBytes))
	}
	return resp, err
}

func HttpDoAuthenticatedDeleteRequest(url string, accessToken string) (*http.Response, error) {
	resp, err := HttpDoAuthenticatedRequest(url, "DELETE", nil, accessToken)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Error API (HTTP %d) : %s", resp.StatusCode, string(bodyBytes))
	}
	return resp, err
}
