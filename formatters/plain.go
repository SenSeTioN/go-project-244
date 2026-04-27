package formatters

import (
	"fmt"
	"strings"

	"code/diff"
)

func formatPlain(nodes []diff.Node) string {
	var lines []string
	collectPlain(&lines, nodes, "")
	return strings.Join(lines, "\n")
}

func collectPlain(lines *[]string, nodes []diff.Node, prefix string) {
	for _, n := range nodes {
		path := n.Key
		if prefix != "" {
			path = prefix + "." + n.Key
		}

		switch n.Status {
		case diff.Added:
			*lines = append(*lines, fmt.Sprintf("Property '%s' was added with value: %s", path, plainValue(n.Value)))
		case diff.Removed:
			*lines = append(*lines, fmt.Sprintf("Property '%s' was removed", path))
		case diff.Changed:
			*lines = append(*lines, fmt.Sprintf("Property '%s' was updated. From %s to %s", path, plainValue(n.OldValue), plainValue(n.NewValue)))
		case diff.Nested:
			collectPlain(lines, n.Children, path)
		case diff.Unchanged:
			// skip
		}
	}
}

func plainValue(v any) string {
	switch val := v.(type) {
	case map[string]any:
		return "[complex value]"
	case string:
		return fmt.Sprintf("'%s'", val)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", val)
	}
}
