package thing

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const thingFileName = "./thingiverse.yml"

var thingFieldOrder = []string{
	"id",
	"name",
	"category",
	"license",
	"is_wip",
	"tags",
	"image_files",
	"model_files",
	"instructions",
	"description",
}

func (tp *Thing) Save() error {
	original, err := os.ReadFile(thingFileName)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var out string
	if err != nil || len(bytes.TrimSpace(original)) == 0 {
		out = renderThingFile(tp)
	} else {
		out, err = patchThingFile(original, tp)
		if err != nil {
			return err
		}
	}

	if !strings.HasSuffix(out, "\n") {
		out += "\n"
	}

	return os.WriteFile(thingFileName, []byte(out), 0644)
}

func patchThingFile(original []byte, tp *Thing) (string, error) {
	var existing Thing
	if err := yaml.Unmarshal(original, &existing); err != nil {
		return "", err
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(original, &doc); err != nil {
		return "", err
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return "", fmt.Errorf("%s root is not a mapping", thingFileName)
	}

	replacements := map[string]string{}
	for _, key := range changedKeys(&existing, tp) {
		replacements[key] = renderField(key, tp)
	}

	missing := missingKeys(doc.Content[0])
	for _, key := range missing {
		if _, already := replacements[key]; already {
			continue
		}
		if isZeroField(key, tp) {
			continue
		}
		replacements[key] = renderField(key, tp)
	}

	if len(replacements) == 0 {
		return string(original), nil
	}

	return applyReplacements(string(original), doc.Content[0], replacements)
}

func changedKeys(old, new *Thing) []string {
	var keys []string
	if old.Id != new.Id {
		keys = append(keys, "id")
	}
	if old.Name != new.Name {
		keys = append(keys, "name")
	}
	if old.Category != new.Category {
		keys = append(keys, "category")
	}
	if old.License != new.License {
		keys = append(keys, "license")
	}
	if old.IsWip != new.IsWip {
		keys = append(keys, "is_wip")
	}
	if !stringSlicesEqual(old.Tags, new.Tags) {
		keys = append(keys, "tags")
	}
	if !filesEqual(old.ImageFiles, new.ImageFiles) {
		keys = append(keys, "image_files")
	}
	if !filesEqual(old.ModelFiles, new.ModelFiles) {
		keys = append(keys, "model_files")
	}
	if old.Instructions != new.Instructions {
		keys = append(keys, "instructions")
	}
	if old.Description != new.Description {
		keys = append(keys, "description")
	}
	return keys
}

func missingKeys(mapping *yaml.Node) []string {
	present := map[string]bool{}
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		present[mapping.Content[i].Value] = true
	}
	var keys []string
	for _, key := range thingFieldOrder {
		if !present[key] {
			keys = append(keys, key)
		}
	}
	return keys
}

func isZeroField(key string, tp *Thing) bool {
	switch key {
	case "id", "category":
		return fieldInt(key, tp) == 0
	case "is_wip":
		return !tp.IsWip
	case "tags":
		return len(tp.Tags) == 0
	case "image_files":
		return len(tp.ImageFiles) == 0
	case "model_files":
		return len(tp.ModelFiles) == 0
	default:
		return fieldString(key, tp) == ""
	}
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func filesEqual(a, b []ThingFile) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].LocalPath != b[i].LocalPath {
			return false
		}
	}
	return true
}

type keySpan struct {
	name  string
	start int
	end   int
}

func topLevelSpans(mapping *yaml.Node, lineCount int) []keySpan {
	var keys []keySpan
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		k := mapping.Content[i]
		keys = append(keys, keySpan{name: k.Value, start: k.Line})
	}
	for i := range keys {
		if i+1 < len(keys) {
			keys[i].end = keys[i+1].start - 1
		} else {
			keys[i].end = lineCount
		}
	}
	return keys
}

func applyReplacements(original string, mapping *yaml.Node, replacements map[string]string) (string, error) {
	lines := strings.Split(strings.ReplaceAll(original, "\r\n", "\n"), "\n")
	spans := topLevelSpans(mapping, len(lines))

	type op struct {
		start int
		end   int
		text  string
	}
	var ops []op
	var appends []string

	for key, rendered := range replacements {
		found := false
		for _, span := range spans {
			if span.name != key {
				continue
			}
			end := span.end
			for end >= span.start && strings.TrimSpace(lines[end-1]) == "" {
				end--
			}
			ops = append(ops, op{start: span.start, end: end, text: rendered})
			found = true
			break
		}
		if !found {
			appends = append(appends, rendered)
		}
	}

	sort.Slice(ops, func(i, j int) bool { return ops[i].start < ops[j].start })

	var out []string
	cursor := 1
	for _, op := range ops {
		if op.start < 1 || op.end > len(lines) || op.start > op.end {
			return "", fmt.Errorf("invalid YAML span for rewrite")
		}
		out = append(out, lines[cursor-1:op.start-1]...)
		out = append(out, strings.Split(op.text, "\n")...)
		cursor = op.end + 1
	}
	if cursor <= len(lines) {
		out = append(out, lines[cursor-1:]...)
	}

	result := strings.Join(out, "\n")
	if len(appends) > 0 {
		sort.Slice(appends, func(i, j int) bool {
			return fieldOrderIndex(appends[i]) < fieldOrderIndex(appends[j])
		})
		if !strings.HasSuffix(result, "\n") {
			result += "\n"
		}
		result += strings.Join(appends, "\n") + "\n"
	}
	return result, nil
}

func fieldOrderIndex(rendered string) int {
	key := rendered
	if i := strings.IndexByte(rendered, ':'); i >= 0 {
		key = rendered[:i]
	}
	for i, name := range thingFieldOrder {
		if name == key {
			return i
		}
	}
	return len(thingFieldOrder)
}

func renderThingFile(tp *Thing) string {
	parts := make([]string, 0, len(thingFieldOrder))
	for _, key := range thingFieldOrder {
		parts = append(parts, renderField(key, tp))
	}
	return strings.Join(parts, "\n") + "\n"
}

func renderField(key string, tp *Thing) string {
	switch key {
	case "id", "category":
		return fmt.Sprintf("%s: %d", key, fieldInt(key, tp))
	case "is_wip":
		return fmt.Sprintf("is_wip: %t", tp.IsWip)
	case "tags":
		return renderStringList("tags", tp.Tags)
	case "image_files":
		return renderFileList("image_files", tp.ImageFiles)
	case "model_files":
		return renderFileList("model_files", tp.ModelFiles)
	default:
		return renderStringField(key, fieldString(key, tp))
	}
}

func fieldInt(key string, tp *Thing) int {
	if key == "category" {
		return tp.Category
	}
	return tp.Id
}

func fieldString(key string, tp *Thing) string {
	switch key {
	case "name":
		return tp.Name
	case "license":
		return tp.License
	case "instructions":
		return tp.Instructions
	case "description":
		return tp.Description
	default:
		return ""
	}
}

func renderStringField(key, value string) string {
	if value == "" {
		return key + `: ""`
	}
	if strings.Contains(value, "\n") {
		return renderLiteralField(key, value)
	}
	if needsQuote(value) {
		return key + ": " + strconv.Quote(value)
	}
	return key + ": " + value
}

func renderLiteralField(key, value string) string {
	chomp := "|"
	body := value
	if strings.HasSuffix(value, "\n") {
		body = strings.TrimSuffix(value, "\n")
	} else {
		chomp = "|-"
	}

	var b strings.Builder
	b.WriteString(key)
	b.WriteString(": ")
	b.WriteString(chomp)
	b.WriteByte('\n')
	for _, line := range strings.Split(body, "\n") {
		b.WriteString("  ")
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func renderStringList(key string, items []string) string {
	if len(items) == 0 {
		return key + ": []"
	}
	var b strings.Builder
	b.WriteString(key)
	b.WriteString(":\n")
	for _, item := range items {
		b.WriteString("  - ")
		if needsQuote(item) {
			b.WriteString(strconv.Quote(item))
		} else {
			b.WriteString(item)
		}
		b.WriteByte('\n')
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func renderFileList(key string, files []ThingFile) string {
	if len(files) == 0 {
		return key + ": []"
	}
	var b strings.Builder
	b.WriteString(key)
	b.WriteString(":\n")
	for _, file := range files {
		b.WriteString("  - local_path: ")
		if needsQuote(file.LocalPath) {
			b.WriteString(strconv.Quote(file.LocalPath))
		} else {
			b.WriteString(file.LocalPath)
		}
		b.WriteByte('\n')
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func needsQuote(s string) bool {
	if s == "" || s == "true" || s == "false" || s == "null" || s == "~" {
		return true
	}
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	if strings.TrimSpace(s) != s {
		return true
	}
	return strings.ContainsAny(s, ":#{}[],&*!|>'\"%@`")
}
