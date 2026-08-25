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

	if strings.TrimSpace(tp.Name) == "" {
		return fmt.Errorf("Name is required")
	}

	if tp.Category == 0 {
		return fmt.Errorf("Category is required")
	}

	if strings.TrimSpace(tp.License) == "" {
		return fmt.Errorf("License is required")
	}

	if strings.TrimSpace(tp.Description) == "" {
		return fmt.Errorf("Description is required")
	}

	return nil
}
