package thing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"

	"github.com/Akrobate/thingiverse-cli/pkg/utils"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

const apiBaseURL = "https://api.thingiverse.com"

type ThingResponse struct {
	ID int `json:"id"`
}

type ThingFile struct {
	LocalPath string `json:"local_path" yaml:"local_path"`
	LocalHash string
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
	ImageFiles   []ThingFile `json:"image_files" yaml:"image_files"`
	ModelFiles   []ThingFile `json:"model_files" yaml:"model_files"`
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
	err = yaml.Unmarshal(data, tp)
	if err != nil {
		return err
	}
	return nil
}

func (tp *Thing) GenerateHashFiles() error {
	for index, _ := range tp.ModelFiles {
		hash, err := utils.CalculateFileHash(tp.ModelFiles[index].LocalPath)
		if err != nil {
			return err
		}
		tp.ModelFiles[index].LocalHash = hash
	}
	return nil
}

func (tp *Thing) CheckFilesExists() error {
	if err := tp.Load(); err != nil {
		return fmt.Errorf("Cannot load thingiverse.yml file in current folder \n%w", err)
	}

	fmt.Println("Image files")
	for _, item := range tp.ImageFiles {
		if utils.FileExists(item.LocalPath) {
			fmt.Printf("[OK]\t%s\n", item.LocalPath)
		} else {
			fmt.Printf("[ERROR]\t%s\n", item.LocalPath)
		}
	}
	fmt.Println("Model files")
	for _, item := range tp.ModelFiles {
		if utils.FileExists(item.LocalPath) {
			fmt.Printf("[OK]\t%s\n", item.LocalPath)
		} else {
			fmt.Printf("[ERROR]\t%s\n", item.LocalPath)
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

	if err := tp.DeleteAllFilesAndImages(accessToken); err != nil {
		return fmt.Errorf("Error updating thing\n%w", err)
	}

	if err := tp.UploadAllFilesAndImages(accessToken); err != nil {
		return fmt.Errorf("Error updating thing\n%w", err)
	}

	return nil
}

func (tp *Thing) CompareAndUpdateFiles(accessToken string) error {

	if err := tp.Load(); err != nil {
		return fmt.Errorf("Cannot load thingiverse.yml file in current folder \n%w", err)
	}
	if err := tp.GenerateHashFiles(); err != nil {
		return fmt.Errorf("Cannot GenerateHashFiles\n%w", err)
	}

	apiFiles, err := GetFilesAPI(tp.Id, accessToken)
	if err != nil {
		return fmt.Errorf("Cannot GetFilesAPI \n%w", err)
	}

	allUniqFilesNamesList := lo.Uniq(append(
		lo.Map(tp.ModelFiles, func(img ThingFile, _ int) string {
			return filepath.Base(img.LocalPath)
		}),
		lo.Map(*apiFiles, func(img FileGetResponse, _ int) string {
			return img.Name
		})...,
	))

	type MixedReponseFile struct {
		LocalPath       string
		FileName        string
		ApiId           int
		ApiImageId      int
		HashMatch       bool
		FoundInApi      bool
		FoundInLocal    bool
		ToDeleteOnApi   bool
		ToCreateOnApi   bool
		ToReuploadOnApi bool
	}

	mixedCommontFilesResponse := lo.Map(allUniqFilesNamesList, func(filename string, _ int) MixedReponseFile {
		apiItem, foundInApi := lo.Find(*apiFiles, func(item FileGetResponse) bool {
			return item.Name == filename
		})

		modelItem, foundInLocal := lo.Find(tp.ModelFiles, func(item ThingFile) bool {
			return filepath.Base(item.LocalPath) == filename
		})

		hashMatch := apiItem.Hash == modelItem.LocalHash

		var resp = MixedReponseFile{
			LocalPath:       modelItem.LocalPath,
			FileName:        apiItem.Name,
			ApiId:           apiItem.Id,
			HashMatch:       hashMatch,
			FoundInApi:      foundInApi,
			FoundInLocal:    foundInLocal,
			ToDeleteOnApi:   foundInApi && !foundInLocal,
			ToCreateOnApi:   !foundInApi && foundInLocal,
			ToReuploadOnApi: foundInApi && foundInLocal && !hashMatch,
		}

		if foundInApi {
			resp.ApiImageId = apiItem.DefaultImage.Id
		}

		return resp
	})

	fmt.Println("--------------- Mixed -----------------------")
	for _, item := range mixedCommontFilesResponse {
		fmt.Printf("ApiId %d \t HashMatch: %t \t FoundInApi: %t \t FoundInLocal: %t \t %s\t %s \t ApiImageId: %d\n",
			item.ApiId, item.HashMatch, item.FoundInApi, item.FoundInLocal, item.LocalPath, item.FileName, item.ApiImageId)
	}

	filterFuncHighlevel := func(prop string) []MixedReponseFile {
		return lo.Filter(mixedCommontFilesResponse, func(item MixedReponseFile, _ int) bool {
			v := reflect.ValueOf(item)
			field := v.FieldByName(prop)
			return field.Bool()
		})
	}

	fmt.Println("")
	fmt.Println("--------------- ToDeleteOnApi -----------------------")
	for _, item := range filterFuncHighlevel("ToDeleteOnApi") {
		fmt.Printf("ApiId %d \t HashMatch: %t \t FoundInApi: %t \t FoundInLocal: %t \t %s\t %s \t ApiImageId: %d\n",
			item.ApiId, item.HashMatch, item.FoundInApi, item.FoundInLocal, item.LocalPath, item.FileName, item.ApiImageId)
	}

	fmt.Println("")
	fmt.Println("--------------- ToCreateOnApi -----------------------")
	for _, item := range filterFuncHighlevel("ToCreateOnApi") {
		fmt.Printf("ApiId %d \t HashMatch: %t \t FoundInApi: %t \t FoundInLocal: %t \t %s\t %s \t ApiImageId: %d\n",
			item.ApiId, item.HashMatch, item.FoundInApi, item.FoundInLocal, item.LocalPath, item.FileName, item.ApiImageId)
	}

	fmt.Println("")
	fmt.Println("--------------- ToReuploadOnApi -----------------------")
	for _, item := range filterFuncHighlevel("ToReuploadOnApi") {
		fmt.Printf("ApiId %d \t HashMatch: %t \t FoundInApi: %t \t FoundInLocal: %t \t %s\t %s \t ApiImageId: %d\n",
			item.ApiId, item.HashMatch, item.FoundInApi, item.FoundInLocal, item.LocalPath, item.FileName, item.ApiImageId)
	}

	return nil
}

func (tp *Thing) UploadAllFilesAndImages(accessToken string) error {
	if err := tp.Load(); err != nil {
		return fmt.Errorf("Cannot load thingiverse.yml file in current folder \n%w", err)
	}
	fmt.Println("Upload image files / model files")
	for _, item := range append(tp.ImageFiles, tp.ModelFiles...) {
		if utils.FileExists(item.LocalPath) {
			err := UploadFileProcess(tp.Id, item.LocalPath, accessToken)
			if err != nil {
				return err
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

// @todo: recode
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
