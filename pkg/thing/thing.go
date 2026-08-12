package thing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
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
	LocalHash string `json:"-" yaml:"-"`
}

type Thing struct {
	Id           int         `json:"-" yaml:"id"`
	Name         string      `json:"name" yaml:"name"`
	Category     int         `json:"category" yaml:"category"`
	License      string      `json:"license" yaml:"license"`
	IsWip        bool        `json:"is_wip" yaml:"is_wip"`
	Tags         []string    `json:"tags" yaml:"tags"`
	ImageFiles   []ThingFile `json:"image_files" yaml:"image_files"`
	ModelFiles   []ThingFile `json:"model_files" yaml:"model_files"`
	Instructions string      `json:"instructions" yaml:"instructions"`
	Description  string      `json:"description" yaml:"description"`
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

func (tp *Thing) CompareAndUpdateFiles(accessToken string, dryRun bool, debug bool) error {

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

	debugPrint := func(item MixedReponseFile) {
		fmt.Printf("ApiId %d \t HashMatch: %t \t FoundInApi: %t \t FoundInLocal: %t \t %s\t %s \t ApiImageId: %d\n",
			item.ApiId, item.HashMatch, item.FoundInApi, item.FoundInLocal, item.LocalPath, item.FileName, item.ApiImageId)
	}

	if debug {
		for _, item := range mixedCommontFilesResponse {
			debugPrint(item)
		}
	}

	filterFuncHighlevel := func(prop string) []MixedReponseFile {
		return lo.Filter(mixedCommontFilesResponse, func(item MixedReponseFile, _ int) bool {
			v := reflect.ValueOf(item)
			field := v.FieldByName(prop)
			return field.Bool()
		})
	}

	fmt.Print("Deleting on Thingiverse...\n")
	for _, item := range filterFuncHighlevel("ToDeleteOnApi") {
		if debug {
			debugPrint(item)
		}

		if !dryRun {
			if err := DeleteImageAPI(item.ApiImageId, tp.Id, accessToken); err != nil {
				fmt.Printf("[Error] DeleteImageAPI %v\n", err)
			}
		}

		if !dryRun {
			if err := DeleteFileAPI(item.ApiId, tp.Id, accessToken); err != nil {
				fmt.Printf("[Error] DeleteFileAPI %v\n", err)
			}
		}

		fmt.Printf("[OK] %s\n", item.FileName)
	}

	fmt.Print("Creating on Thingiverse...\n")
	for _, item := range filterFuncHighlevel("ToCreateOnApi") {
		if debug {
			debugPrint(item)
		}

		if !dryRun {
			if err := UploadFileProcess(tp.Id, item.LocalPath, accessToken); err != nil {
				fmt.Printf("[Error] CreateFileAPI %v\n", err)
			}
		}

		fmt.Printf("[OK] %s\n", filepath.Base(item.LocalPath))
	}

	fmt.Print("Reupload on Thingiverse...\n")
	for _, item := range filterFuncHighlevel("ToReuploadOnApi") {
		if debug {
			debugPrint(item)
		}

		if !dryRun {
			if err := DeleteImageAPI(item.ApiImageId, tp.Id, accessToken); err != nil {
				fmt.Printf("[Error] DeleteImageAPI %v\n", err)
			}
		}

		if !dryRun {
			if err := DeleteFileAPI(item.ApiId, tp.Id, accessToken); err != nil {
				fmt.Printf("[Error] DeleteFileAPI %v\n", err)
			}
		}

		if !dryRun {
			if err := UploadFileProcess(tp.Id, item.LocalPath, accessToken); err != nil {
				fmt.Printf("[Error] CreateFileAPI %v\n", err)
			}
		}

		fmt.Printf("[OK] %s\n", filepath.Base(item.LocalPath))
	}

	return nil
}

func (tp *Thing) DeleteAndUpdateAllImages(accessToken string, dryRun bool, debug bool) error {

	fmt.Println("Deleting images")

	// @todo: Check if the images returns also the rendered stl files. Should not be totaly deleted
	images, err := GetImagesAPI(tp.Id, accessToken)
	if err != nil {
		return err
	}

	for _, item := range *images {
		err := DeleteImageAPI(item.Id, tp.Id, accessToken)
		if err != nil {
			fmt.Printf("[ERROR] %d\t %s %w\n", item.Id, item.Name, err)
		} else {
			fmt.Printf("[OK] %d\t %s\n", item.Id, item.Name)
		}
	}

	for _, item := range tp.ImageFiles {
		if utils.FileExists(item.LocalPath) {
			err := UploadFileProcess(tp.Id, item.LocalPath, accessToken)
			if err != nil {
				return err
			}
		}
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

// @todo, To test
func (tp *Thing) UpdateOrderFilesAndImage(id int, accessToken string) error {

	type ImageFileOrderItem struct {
		Id   int `json:"id"`
		Rank int `json:"rank"`
	}

	type ThingFileImageOrderUpdateRequest struct {
		Files  []ImageFileOrderItem `json:"files"`
		Images []ImageFileOrderItem `json:"images"`
	}

	var request ThingFileImageOrderUpdateRequest

	images, err := GetGalleriesFilesWithoutModelsPreviews(tp.Id, accessToken)
	files, err := GetFilesAPI(tp.Id, accessToken)

	if err != nil {
		return fmt.Errorf("Error GetGalleriesFilesWithoutModelsPreviews %w", err)
	}

	indexImage := 0

	for _, local_item := range tp.ImageFiles {
		foundApiItem, foundInApi := lo.Find(*images, func(apiItem ImageGetResponse) bool {
			return apiItem.Name == filepath.Base(local_item.LocalPath)
		})

		if !foundInApi {
			return fmt.Errorf("Error Not found on api %s", filepath.Base(local_item.LocalPath))
		}

		var orderImageItem = ImageFileOrderItem{
			Id:   foundApiItem.Id,
			Rank: indexImage,
		}
		request.Images = append(request.Images, orderImageItem)
		indexImage++
	}

	indexFile := 0
	for _, local_item := range tp.ModelFiles {
		foundApiItem, foundInApi := lo.Find(*files, func(apiItem FileGetResponse) bool {
			return apiItem.Name == filepath.Base(local_item.LocalPath)
		})

		if !foundInApi {
			return fmt.Errorf("Error Not found on api %s", filepath.Base(local_item.LocalPath))
		}

		var orderFileItem = ImageFileOrderItem{
			Id:   foundApiItem.Id,
			Rank: indexFile,
		}
		request.Files = append(request.Files, orderFileItem)
		indexFile++

		var orderImageItem = ImageFileOrderItem{
			Id:   foundApiItem.DefaultImage.Id,
			Rank: indexImage,
		}
		request.Images = append(request.Images, orderImageItem)

		indexImage++
	}

	if err := UpdateAPI(tp.Id, &request, accessToken); err != nil {
		return fmt.Errorf("Error updating thing\n%w", err)
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

func (tp *Thing) AutosetFilesAndImages(rootDir string, mode string) error {

	// var list []string

	if mode == "images" {
		tp.ImageFiles = nil
		err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() || filepath.Ext(path) != ".png" {
				return nil
			}

			var tf = ThingFile{
				LocalPath: path,
			}
			tp.ImageFiles = append(tp.ImageFiles, tf)

			return nil
		})

		if err != nil {
			fmt.Println("Erreur :", err)
		}
	}

	if mode == "files" {
		tp.ModelFiles = nil
		err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() || filepath.Ext(path) != ".stl" {
				return nil
			}

			var tf = ThingFile{
				LocalPath: path,
			}
			tp.ModelFiles = append(tp.ModelFiles, tf)
			return nil
		})

		if err != nil {
			fmt.Println("Erreur :", err)
		}
	}

	// for _, item := range list {
	// 	fmt.Println(item)
	// }

	tp.Save()

	return nil
}
