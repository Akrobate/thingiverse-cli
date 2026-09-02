package thing

import (
	"os"
	"strings"
	"testing"
)

const sampleThingYAML = `id: 7387159
name: TP4056 LiPo Charger Panel / Surface Mount Bracket
category: 66
license: cc
is_wip: false
tags:
  - Tp4056
  - tp4056_case
image_files:
  - local_path: ./assets/photos/printed_preview_2.jpg
  - local_path: ./assets/photos/printed_preview_1.jpg
model_files:
  - local_path: ./opm_stl_files/tp4056FixationPiece.stl
instructions: ""
description: |
  This is a modular, 2-piece mounting bracket designed to securely hold a standard **TP4056 LiPo battery charging module** in place.

  ### Features

  * **Mounting Style:** Panel mount
  * **Hardware Required:**
    * 2x M3 screws
`

func chdirTemp(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(old)
	})
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
}

func TestSavePreservesMultilineDescriptionAfterAutoset(t *testing.T) {
	chdirTemp(t)

	if err := os.WriteFile(thingFileName, []byte(sampleThingYAML), 0644); err != nil {
		t.Fatal(err)
	}

	tp, err := NewThing()
	if err != nil {
		t.Fatal(err)
	}
	if err := tp.Load(); err != nil {
		t.Fatal(err)
	}

	tp.ImageFiles = append(tp.ImageFiles, ThingFile{LocalPath: "./opm_png_files/example.png"})
	if err := tp.Save(); err != nil {
		t.Fatal(err)
	}

	out, err := os.ReadFile(thingFileName)
	if err != nil {
		t.Fatal(err)
	}
	got := string(out)

	if strings.HasPrefix(got, "---") {
		t.Fatalf("rewritten file should not start with document marker, got:\n%s", got)
	}
	if !strings.Contains(got, "description: |") {
		t.Fatalf("description should stay a literal block scalar, got:\n%s", got)
	}
	if !strings.Contains(got, "instructions: \"\"") {
		t.Fatalf("empty instructions should stay a quoted empty string, got:\n%s", got)
	}
	if strings.Contains(got, "instructions: |") {
		t.Fatalf("empty instructions should not become a literal block, got:\n%s", got)
	}
	if !strings.Contains(got, "  This is a modular, 2-piece mounting bracket") {
		t.Fatalf("description body was not preserved, got:\n%s", got)
	}
	if !strings.Contains(got, "  ### Features") {
		t.Fatalf("description markdown headings were not preserved, got:\n%s", got)
	}
	if !strings.Contains(got, "  - local_path: ./opm_png_files/example.png") {
		t.Fatalf("new image_files entry should use block list style, got:\n%s", got)
	}
	if strings.Contains(got, "image_files: [{") {
		t.Fatalf("image_files should not be flow style, got:\n%s", got)
	}

	reloaded, err := NewThing()
	if err != nil {
		t.Fatal(err)
	}
	if err := reloaded.Load(); err != nil {
		t.Fatal(err)
	}
	if reloaded.Description != tp.Description {
		t.Fatalf("description content changed\nwant: %q\ngot:  %q", tp.Description, reloaded.Description)
	}
	if got, want := len(reloaded.ImageFiles), 3; got != want {
		t.Fatalf("image_files count: got %d want %d", got, want)
	}
}

func TestSaveWritesMultilineDescriptionOnNewFile(t *testing.T) {
	chdirTemp(t)

	tp := &Thing{
		Id:          1,
		Name:        "Example",
		Category:    66,
		License:     "cc",
		Description: "Line one\n\nLine two\n",
	}
	if err := tp.Save(); err != nil {
		t.Fatal(err)
	}

	out, err := os.ReadFile(thingFileName)
	if err != nil {
		t.Fatal(err)
	}
	got := string(out)

	if !strings.Contains(got, "description: |") {
		t.Fatalf("new multiline description should use literal style, got:\n%s", got)
	}
	if !strings.Contains(got, "instructions: \"\"") {
		t.Fatalf("empty instructions should be a quoted empty string, got:\n%s", got)
	}
}

func TestSavePreservesOriginalDescriptionBytes(t *testing.T) {
	chdirTemp(t)

	src := `id: 1
image_files:
  - local_path: ./old.png
instructions: ""
description: |
  Line with trailing space. 
  Emoji ⚙️ and fences:

  ` + "```bash\n" + `  opm install
  ` + "```\n"

	if err := os.WriteFile(thingFileName, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	tp, err := NewThing()
	if err != nil {
		t.Fatal(err)
	}
	if err := tp.Load(); err != nil {
		t.Fatal(err)
	}
	tp.ImageFiles = []ThingFile{{LocalPath: "./new.png"}}
	if err := tp.Save(); err != nil {
		t.Fatal(err)
	}

	out, err := os.ReadFile(thingFileName)
	if err != nil {
		t.Fatal(err)
	}
	got := string(out)
	if !strings.Contains(got, "description: |") {
		t.Fatalf("description should stay a literal block, got:\n%s", got)
	}
	if strings.Contains(got, `description: "`) {
		t.Fatalf("description should not become a quoted string, got:\n%s", got)
	}
	if !strings.Contains(got, "  Line with trailing space. ") {
		t.Fatalf("trailing spaces in description should be kept, got:\n%s", got)
	}
	if !strings.Contains(got, "  - local_path: ./new.png") {
		t.Fatalf("image_files was not updated, got:\n%s", out)
	}
	if strings.Contains(got, "name: \"\"") || strings.Contains(got, "category: 0") {
		t.Fatalf("zero-value fields missing from the file should not be appended, got:\n%s", got)
	}

	reloaded, err := NewThing()
	if err != nil {
		t.Fatal(err)
	}
	if err := reloaded.Load(); err != nil {
		t.Fatal(err)
	}
	if reloaded.Description != tp.Description {
		t.Fatalf("description content changed\nwant: %q\ngot:  %q", tp.Description, reloaded.Description)
	}
}

func TestSaveRewritesDescriptionAsLiteralBlock(t *testing.T) {
	chdirTemp(t)

	if err := os.WriteFile(thingFileName, []byte("id: 1\ndescription: old\n"), 0644); err != nil {
		t.Fatal(err)
	}

	tp := &Thing{Id: 1, Description: "Line with trailing space. \nEmoji ⚙️\n"}
	if err := tp.Save(); err != nil {
		t.Fatal(err)
	}

	out, err := os.ReadFile(thingFileName)
	if err != nil {
		t.Fatal(err)
	}
	got := string(out)
	if !strings.Contains(got, "description: |") {
		t.Fatalf("updated description should use literal style, got:\n%s", got)
	}
	if strings.Contains(got, `description: "`) {
		t.Fatalf("updated description should not be double-quoted, got:\n%s", got)
	}
	if !strings.Contains(got, "  Line with trailing space. ") {
		t.Fatalf("trailing spaces in description should be kept, got:\n%s", got)
	}
}

func TestMarshalYAMLDoesNotDependOnThingFields(t *testing.T) {
	type sample struct {
		Title string   `yaml:"title"`
		Notes string   `yaml:"notes"`
		Items []string `yaml:"items"`
	}

	out, err := marshalYAML(&sample{
		Title: "demo",
		Notes: "line 1\n\nline 2\n",
		Items: []string{"one", "two"},
	})
	if err != nil {
		t.Fatal(err)
	}
	got := string(out)
	if !strings.Contains(got, "notes: |") {
		t.Fatalf("generic marshal should use literal blocks, got:\n%s", got)
	}
	if !strings.Contains(got, "  - one") || !strings.Contains(got, "  - two") {
		t.Fatalf("generic marshal should use block lists, got:\n%s", got)
	}
}

