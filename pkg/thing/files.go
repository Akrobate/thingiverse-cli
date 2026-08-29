package thing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/Akrobate/thingiverse-cli/pkg/utils"
)

type pendingUploadResponse struct {
	ID int `json:"id"`
}

type pendingUploadItem struct {
	ID   int `json:"id"`
	Rank int `json:"rank"`
}

type finalizeFilesRequest struct {
	PendingUploads []pendingUploadItem `json:"pending_uploads"`
	TargetID       int                 `json:"target_id"`
	TargetType     string              `json:"target_type"`
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
	pendingID, err := UploadPendingFileAPI(localPath, accessToken)
	if err != nil {
		return fmt.Errorf("UploadPendingFileAPI error\n%w", err)
	}

	if err := FinalizePendingFilesAPI(thingID, []int{pendingID}, accessToken); err != nil {
		return fmt.Errorf("FinalizePendingFilesAPI error\n%w", err)
	}
	return nil
}

func UploadPendingFileAPI(localPath string, accessToken string) (int, error) {
	file, err := os.Open(localPath)
	if err != nil {
		return 0, fmt.Errorf("cannot open file %s: %w", localPath, err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(localPath))
	if err != nil {
		return 0, fmt.Errorf("failed creating multipart: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return 0, fmt.Errorf("failed copy file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return 0, fmt.Errorf("failed closing multipart writer: %w", err)
	}

	url := fmt.Sprintf("%s/files/0/uploadFile", apiBaseURL)
	resp, err := utils.HttpDoAuthenticatedRequestWithContentType(
		url,
		"POST",
		body,
		writer.FormDataContentType(),
		accessToken,
	)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed reading upload response: %w", err)
	}

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("API uploadFile error (HTTP %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var pending pendingUploadResponse
	if err := json.Unmarshal(bodyBytes, &pending); err != nil {
		return 0, fmt.Errorf("parse uploadFile response: %w", err)
	}
	if pending.ID == 0 {
		return 0, fmt.Errorf("uploadFile response missing id: %s", string(bodyBytes))
	}
	return pending.ID, nil
}

func FinalizePendingFilesAPI(thingID int, pendingIDs []int, accessToken string) error {
	items := make([]pendingUploadItem, 0, len(pendingIDs))
	for i, id := range pendingIDs {
		items = append(items, pendingUploadItem{ID: id, Rank: i})
	}

	jsonData, err := json.Marshal(finalizeFilesRequest{
		PendingUploads: items,
		TargetID:       thingID,
		TargetType:     "thing",
	})
	if err != nil {
		return fmt.Errorf("error JSON serialize: %w", err)
	}

	url := fmt.Sprintf("%s/files/0/FinalizeFiles", apiBaseURL)
	resp, err := utils.HttpDoAuthenticatedPostRequest(url, jsonData, accessToken)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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
