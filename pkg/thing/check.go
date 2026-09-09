package thing

import (
	"fmt"
	"strings"
)

func (tp *Thing) CheckThingParams() error {

	// 	Id           int         `json:"-" yaml:"id"`
	// 	Name         string      `json:"name" yaml:"name"`
	// 	Category     int         `json:"category" yaml:"category"`
	// 	License      string      `json:"license" yaml:"license"`
	// 	IsWip        bool        `json:"is_wip" yaml:"is_wip"`
	// 	Tags         []string    `json:"tags" yaml:"tags"`
	// 	ImageFiles   []ThingFile `json:"image_files" yaml:"image_files"`
	// 	ModelFiles   []ThingFile `json:"model_files" yaml:"model_files"`
	// 	Instructions string      `json:"instructions" yaml:"instructions"`
	// 	Description  string      `json:"description" yaml:"description"`

	errMessage := ""
	hadError := false

	if strings.TrimSpace(tp.Name) == "" {
		hadError = true
		errMessage = errMessage + fmt.Sprintf("Name is required")
	}

	if tp.Category == 0 {
		hadError = true
		errMessage = errMessage + fmt.Sprintf("Category is required")
	}

	if strings.TrimSpace(tp.License) == "" {
		hadError = true
		errMessage = errMessage + fmt.Sprintf("License is required")
	}

	if strings.TrimSpace(tp.Description) == "" {
		hadError = true
		errMessage = errMessage + fmt.Sprintf("Description is required")
	}

	if hadError {
		return fmt.Errorf(errMessage)
	}
	return nil
}
