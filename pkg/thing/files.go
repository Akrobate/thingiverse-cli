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

	"github.com/Akrobate/thingiverse-cli/pkg/utils"
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

func UploadFileProcess(thingID int, localPath string, accessToken string) error {
	filename := filepath.Base(localPath)
	creationResponse, err := CreateFileAPI(thingID, filename, accessToken)
	if err != nil {
		return fmt.Errorf("CreateFileAPI error \n%w", err)
	}

	err = UploadToS3(creationResponse.Action, creationResponse.Fields, localPath)
	if err != nil {
		return fmt.Errorf("UploadToS3 error \n%w", err)
	}

	err = FinaliseFileAPI(creationResponse.Fields.SuccessActionRedirect, creationResponse.Fields, accessToken)
	if err != nil {
		return fmt.Errorf("FinaliseFileAPI error \n%w", err)
	}
	return nil
}

func CreateFileAPI(thingID int, filename string, accessToken string) (*ThingFilePostResponse, error) {

	postRequest := struct {
		Filename string `json:"filename"`
	}{Filename: filename}

	jsonData, err := json.Marshal(postRequest)
	if err != nil {
		return nil, fmt.Errorf("Error JSON serialize : %w", err)
	}

	url := fmt.Sprintf("%s/things/%d/files", apiBaseURL, thingID)
	resp, err := utils.HttpDoAuthenticatedPostRequest(url, jsonData, accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
	resp, err := utils.HttpDoAuthenticatedPostRequest(successActionRedirect, jsonData, accessToken)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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
	resp, err := utils.HttpDoAuthenticatedDeleteRequest(url, accessToken)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func GetImagesAPI(thingID int, accessToken string) (*[]ImageGetResponse, error) {
	url := fmt.Sprintf("%s/things/%d/%s", apiBaseURL, thingID, "images")

	resp, err := utils.HttpDoAuthenticatedGetRequest(url, accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var t []ImageGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("parse response problem: %w", err)
	}

	return &t, nil
}

func GetFilesAPI(thingID int, accessToken string) (*[]FileGetResponse, error) {
	url := fmt.Sprintf("%s/things/%d/%s", apiBaseURL, thingID, "files")

	resp, err := utils.HttpDoAuthenticatedGetRequest(url, accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var t []FileGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("parse response problem: %w", err)
	}

	return &t, nil
}

func GetFilePreivewImages(thingID int, accessToken string) (*[]ImageGetResponse, error) {

	files, err := GetFilesAPI(thingID, accessToken)
	if err != nil {
		return nil, fmt.Errorf("Error GetFilesAPI %w", err)
	}

	var results []ImageGetResponse
	for _, item := range *files {
		results = append(results, item.DefaultImage)
	}
	return &results, nil
}

func GetGalleriesFilesWithoutModelsPreviews(thingID int, accessToken string) (*[]ImageGetResponse, error) {
	images, err := GetImagesAPI(thingID, accessToken)
	if err != nil {
		return nil, fmt.Errorf("Cannot GetImagesAPI %w", err)
	}
	filePreviewImages, err := GetFilePreivewImages(thingID, accessToken)
	if err != nil {
		return nil, fmt.Errorf("Cannot GetFilePreivewImages %w", err)
	}
	excludedIDs := make(map[int]bool, len(*filePreviewImages))
	for _, img := range *filePreviewImages {
		excludedIDs[img.Id] = true
	}
	var filteredImages []ImageGetResponse
	for _, img := range *images {
		if !excludedIDs[img.Id] {
			filteredImages = append(filteredImages, img)
		}
	}
	return &filteredImages, nil
}
