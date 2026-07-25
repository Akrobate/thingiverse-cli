package thing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Akrobate/thingiverse-cli/pkg/utils"
	"gopkg.in/yaml.v3"
)

const apiBaseURL = "https://api.thingiverse.com"

type ThingResponse struct {
	ID int `json:"id"`
}

type ThingFile struct {
	LocalPath string `json:"local_path" yaml:"local_path"`
}

type Thing struct {
	Id           int         `json:"-" yaml:"id"`
	Name         string      `json:"name" yaml:"name"`
	Category     int         `json:"category" yaml:"category"`
	License      string      `json:"license" yaml:"license"`
	IsWip        bool        `json:"is_wip" yaml:"is_wip"`
	Tags         []string    `json:"tags" yaml:"tags"`
	Instructions string      `json:"instructions" yaml:"instructions"`
	Description  string      `json:"description" yaml:"description"`
	Files        []ThingFile `json:"files" yaml:"files"`
}

func NewThing() (*Thing, error) {

	return &Thing{}, nil
}

func (tp *Thing) Save() error {
	data, err := yaml.Marshal(tp)
	if err != nil {
		return err
	}

	return os.WriteFile("./thingiverse.yml", data, 0644)
}

func (tp *Thing) Load() error {
	data, err := os.ReadFile("./thingiverse.yml")
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, tp)
}

func (tp *Thing) CheckFilesExists() error {
	if err := tp.Load(); err != nil {
		return fmt.Errorf("Cannot load thingiverse.yml file in current folder \n%w", err)
	}

	for _, item := range tp.Files {
		if utils.FileExists(item.LocalPath) {
			fmt.Printf("[OK]\t%s\n", item.LocalPath)
		} else {
			fmt.Printf("[ERROR]\t%s\n", item.LocalPath)
		}
	}

	return nil
}

func (tp *Thing) UploadFiles(accessToken string) error {
	if err := tp.Load(); err != nil {
		return fmt.Errorf("Cannot load thingiverse.yml file in current folder \n%w", err)
	}

	for _, item := range tp.Files {
		if utils.FileExists(item.LocalPath) {

			filename := filepath.Base(item.LocalPath)
			creationResponse, err := CreateFileAPI(tp.Id, filename, accessToken)
			if err != nil {
				return fmt.Errorf("CreateFileAPI error \n%w", err)
			}

			err = UploadToS3(creationResponse.Action, creationResponse.Fields, item.LocalPath)
			if err != nil {
				return fmt.Errorf("UploadToS3 error \n%w", err)
			}

			err = FinaliseFileAPI(creationResponse.Fields.SuccessActionRedirect, creationResponse.Fields, accessToken)
			if err != nil {
				return fmt.Errorf("FinaliseFileAPI error \n%w", err)
			}
		}
	}

	return nil
}

func (tp *Thing) DeleteAllFilesAndImages(accessToken string) error {
	if err := tp.Load(); err != nil {
		return fmt.Errorf("Cannot load thingiverse.yml file in current folder \n%w", err)
	}

	fmt.Println("Deleting images")
	images, err := GetImagesAPI(tp.Id, accessToken)
	if err != nil {
		return err
	}
	for _, item := range *images {
		err := DeleteImageAPI(item.Id, tp.Id, accessToken)
		if err != nil {
			fmt.Printf("[ERROR] %d\t %s %w\n", item.Id, item.Name, err)
		} else {
			fmt.Printf("[DELETED] %d\t %s\n", item.Id, item.Name)
		}
	}

	fmt.Println("Deleting files")
	files, err := GetFilesAPI(tp.Id, accessToken)
	if err != nil {
		return err
	}
	for _, item := range *files {
		err := DeleteFileAPI(item.Id, tp.Id, accessToken)
		if err != nil {
			fmt.Printf("[ERROR] %d\t %s %w\n", item.Id, item.Name, err)
		} else {
			fmt.Printf("[DELETED] %d\t %s\n", item.Id, item.Name)
		}
	}
	return nil
}

func (tp *Thing) Update(accessToken string) error {

	if err := tp.Load(); err != nil {
		return fmt.Errorf("Cannot load thingiverse.yml file in current folder \n%w", err)
	}

	updateRequest := ThingUpdateRequest{
		Name:        tp.Name,
		Category:    tp.Category,
		Description: tp.Description,
		IsWip:       tp.IsWip,
		Tags:        tp.Tags,
		License:     tp.License,
	}

	if err := UpdateAPI(tp.Id, &updateRequest, accessToken); err != nil {
		return fmt.Errorf("Error updating thing\n%w", err)
	}

	return nil
}

func (tp *Thing) Create(accessToken string) (int, error) {

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

	tp.Id = thingResp.ID

	return thingResp.ID, nil
}
