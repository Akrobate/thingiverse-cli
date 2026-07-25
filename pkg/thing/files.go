package thing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type ThingFilePostResponse struct {
	Action      string         `json:"action"`
	IsImageFile bool           `json:"is_image_file"`
	Fields      S3UploadPolicy `json:"fields"`
}

type S3UploadPolicy struct {
	AWSAccessKeyID        string `json:"AWSAccessKeyId"`
	Bucket                string `json:"bucket"`
	Key                   string `json:"key"`
	ACL                   string `json:"acl"`
	Policy                string `json:"policy"`
	Signature             string `json:"signature"`
	SuccessActionRedirect string `json:"success_action_redirect"`
	ContentType           string `json:"Content-Type"`
	ContentDisposition    string `json:"Content-Disposition"`
}

type ImageGetResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Rank int    `json:"rank"`
}

type FileGetResponse struct {
	Id           int              `json:"id"`
	Name         string           `json:"name"`
	Size         int64            `json:"size"`
	Hash         string           `json:"hash"`
	DefaultImage ImageGetResponse `json:"default_image"`
}

func CreateFileAPI(id int, filename string, accessToken string) (*ThingFilePostResponse, error) {

	postRequest := struct {
		Filename string `json:"filename"`
	}{
		Filename: filename,
	}

	jsonData, err := json.Marshal(postRequest)
	if err != nil {
		return nil, fmt.Errorf("Error JSON serialize : %w", err)
	}

	url := fmt.Sprintf("%s/things/%d/files", apiBaseURL, id)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("Error creating request : %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request failed : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API Error (HTTP %d) : %s", resp.StatusCode, string(bodyBytes))
	}

	var t ThingFilePostResponse
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("parse response problem: %w", err)
	}
	return &t, nil
}

func FinaliseFileAPI(successActionRedirect string, policy S3UploadPolicy, accessToken string) error {

	jsonData, err := json.Marshal(policy)
	if err != nil {
		return fmt.Errorf("Error JSON serialize : %w", err)
	}

	req, err := http.NewRequest("POST", successActionRedirect, bytes.NewBuffer(jsonData))
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

func UploadToS3(uploadURL string, policy S3UploadPolicy, filePath string) error {

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Cannot open file %s: %w", filePath, err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	formFields := []struct {
		key   string
		value string
	}{
		{"AWSAccessKeyId", policy.AWSAccessKeyID},
		{"bucket", policy.Bucket},
		{"key", policy.Key},
		{"acl", policy.ACL},
		{"policy", policy.Policy},
		{"signature", policy.Signature},
		{"success_action_redirect", policy.SuccessActionRedirect},
		{"Content-Type", policy.ContentType},
		{"Content-Disposition", policy.ContentDisposition},
	}

	for _, field := range formFields {
		if err := writer.WriteField(field.key, field.value); err != nil {
			return fmt.Errorf("Error writing field %s: %w", field.key, err)
		}
	}

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("Failed ceating multipart: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("Failed copy file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("Failed closing multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return fmt.Errorf("Error creating request : %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed sending storage S3/GCP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Storage error (%d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func DeleteImageAPI(imageID int, thingID int, accessToken string) error {
	return genericContentDelete("images", imageID, thingID, accessToken)
}

func DeleteFileAPI(imageID int, thingID int, accessToken string) error {
	return genericContentDelete("files", imageID, thingID, accessToken)
}

func genericContentDelete(item string, imageID int, thingID int, accessToken string) error {

	url := fmt.Sprintf("%s/things/%d/%s/%d", apiBaseURL, thingID, item, imageID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("Error creating request : %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed HTTP DELETE : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Error API (HTTP %d) : %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func GetImagesAPI(thingID int, accessToken string) (*[]ImageGetResponse, error) {
	url := fmt.Sprintf("%s/things/%d/%s", apiBaseURL, thingID, "images")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
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

	var t []ImageGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("parse response problem: %w", err)
	}

	return &t, nil
}

func GetFilesAPI(thingID int, accessToken string) (*[]FileGetResponse, error) {
	url := fmt.Sprintf("%s/things/%d/%s", apiBaseURL, thingID, "files")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
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

	var t []FileGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("parse response problem: %w", err)
	}

	return &t, nil
}
