package thing

import (
	"bytes"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const thingFileName = "./thingiverse.yml"

func (tp *Thing) Save() error {
	node, err := toYAMLNode(tp)
	if err != nil {
		return err
	}
	if previous, err := readYAMLNode(thingFileName); err == nil {
		mergeMapping(previous, node)
		node = previous
	}
	return os.WriteFile(thingFileName, renderYAML(node), 0644)
}

// marshalYAML encodes any value as YAML. Multiline strings use a literal
// block (`|`), empty strings stay `""`, and collections stay in block style.
func marshalYAML(v any) ([]byte, error) {
	node, err := toYAMLNode(v)
	if err != nil {
		return nil, err
	}
	return renderYAML(node), nil
}

func renderYAML(node *yaml.Node) []byte {
	var buf bytes.Buffer
	writeNode(&buf, node, 0)
	if buf.Len() == 0 || buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

func toYAMLNode(v any) (*yaml.Node, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return nil, err
	}
	return parseYAMLNode(data)
}

func readYAMLNode(path string) (*yaml.Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, os.ErrNotExist
	}
	return parseYAMLNode(data)
}

func parseYAMLNode(data []byte) (*yaml.Node, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	if len(doc.Content) == 0 {
		return &yaml.Node{Kind: yaml.MappingNode}, nil
	}
	return doc.Content[0], nil
}

func mergeMapping(dst, src *yaml.Node) {
	if dst.Kind != yaml.MappingNode || src.Kind != yaml.MappingNode {
		*dst = *src
		return
	}

	index := map[string]int{}
	for i := 0; i+1 < len(dst.Content); i += 2 {
		index[dst.Content[i].Value] = i
	}

	for i := 0; i+1 < len(src.Content); i += 2 {
		key := src.Content[i].Value
		if j, ok := index[key]; ok {
			if !nodesEqual(dst.Content[j+1], src.Content[i+1]) {
				dst.Content[j+1] = src.Content[i+1]
			}
			continue
		}
		if isZeroNode(src.Content[i+1]) {
			continue
		}
		dst.Content = append(dst.Content, src.Content[i], src.Content[i+1])
	}
}

func isZeroNode(n *yaml.Node) bool {
	switch n.Kind {
	case yaml.SequenceNode, yaml.MappingNode:
		return len(n.Content) == 0
	case yaml.ScalarNode:
		switch n.Tag {
		case "!!str", "!":
			return n.Value == ""
		case "!!bool":
			return n.Value == "false"
		case "!!int", "!!float":
			return n.Value == "0" || n.Value == "0.0"
		default:
			return n.Value == "" || n.Value == "0" || n.Value == "false" || n.Value == "null" || n.Value == "~"
		}
	default:
		return n.Value == ""
	}
}

func nodesEqual(a, b *yaml.Node) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.Kind != b.Kind || a.Value != b.Value || len(a.Content) != len(b.Content) {
		return false
	}
	for i := range a.Content {
		if !nodesEqual(a.Content[i], b.Content[i]) {
			return false
		}
	}
	return true
}

func writeNode(b *bytes.Buffer, n *yaml.Node, indent int) {
	switch n.Kind {
	case yaml.MappingNode:
		for i := 0; i+1 < len(n.Content); i += 2 {
			writeIndent(b, indent)
			writeScalar(b, n.Content[i], indent)
			b.WriteByte(':')
			writeFieldValue(b, n.Content[i+1], indent)
		}
	case yaml.SequenceNode:
		writeSeq(b, n, indent)
	default:
		writeIndent(b, indent)
		if writeScalar(b, n, indent) {
			b.WriteByte('\n')
		}
	}
}

func writeFieldValue(b *bytes.Buffer, n *yaml.Node, indent int) {
	switch n.Kind {
	case yaml.MappingNode:
		if len(n.Content) == 0 {
			b.WriteString(" {}\n")
			return
		}
		b.WriteByte('\n')
		writeNode(b, n, indent+1)
	case yaml.SequenceNode:
		if len(n.Content) == 0 {
			b.WriteString(" []\n")
			return
		}
		b.WriteByte('\n')
		writeSeq(b, n, indent+1)
	default:
		b.WriteByte(' ')
		if writeScalar(b, n, indent) {
			b.WriteByte('\n')
		}
	}
}

func writeSeq(b *bytes.Buffer, n *yaml.Node, indent int) {
	for _, item := range n.Content {
		writeIndent(b, indent)
		b.WriteByte('-')
		switch item.Kind {
		case yaml.MappingNode:
			if len(item.Content) == 0 {
				b.WriteString(" {}\n")
				continue
			}
			b.WriteByte(' ')
			writeScalar(b, item.Content[0], indent)
			b.WriteByte(':')
			writeFieldValue(b, item.Content[1], indent+1)
			for i := 2; i+1 < len(item.Content); i += 2 {
				writeIndent(b, indent+1)
				writeScalar(b, item.Content[i], indent)
				b.WriteByte(':')
				writeFieldValue(b, item.Content[i+1], indent+1)
			}
		case yaml.SequenceNode:
			if len(item.Content) == 0 {
				b.WriteString(" []\n")
				continue
			}
			b.WriteByte('\n')
			writeSeq(b, item, indent+1)
		default:
			b.WriteByte(' ')
			if writeScalar(b, item, indent) {
				b.WriteByte('\n')
			}
		}
	}
}

func writeScalar(b *bytes.Buffer, n *yaml.Node, indent int) bool {
	if isYAMLString(n) {
		switch {
		case n.Value == "":
			b.WriteString(`""`)
			return true
		case strings.Contains(n.Value, "\n"):
			writeLiteral(b, n.Value, indent)
			return false
		case needsQuote(n.Value):
			b.WriteString(strconv.Quote(n.Value))
			return true
		}
	}
	b.WriteString(n.Value)
	return true
}

func writeLiteral(b *bytes.Buffer, value string, indent int) {
	marker := "|"
	body := value
	if strings.HasSuffix(value, "\n") {
		body = strings.TrimSuffix(value, "\n")
	} else {
		marker = "|-"
	}

	b.WriteString(marker)
	b.WriteByte('\n')
	pad := strings.Repeat("  ", indent+1)
	for _, line := range strings.Split(body, "\n") {
		b.WriteString(pad)
		b.WriteString(line)
		b.WriteByte('\n')
	}
}

func writeIndent(b *bytes.Buffer, indent int) {
	b.WriteString(strings.Repeat("  ", indent))
}

func isYAMLString(n *yaml.Node) bool {
	switch n.Tag {
	case "!!str", "!":
		return true
	case "":
		return n.Style&(yaml.DoubleQuotedStyle|yaml.SingleQuotedStyle|yaml.LiteralStyle|yaml.FoldedStyle) != 0
	default:
		return false
	}
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
