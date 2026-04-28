package formatters

import (
	"fmt"
	"sort"
	"strings"

	"code/internal/diff"
)

const stylishIndent = 4

func formatStylish(nodes []diff.Node) string {
	var b strings.Builder
	b.WriteString("{\n")
	renderStylish(&b, nodes, 1)
	b.WriteString("}")
	return b.String()
}

func renderStylish(b *strings.Builder, nodes []diff.Node, depth int) {
	indent := strings.Repeat(" ", depth*stylishIndent)
	markerIndent := strings.Repeat(" ", depth*stylishIndent-2)

	for _, n := range nodes {
		switch n.Status {
		case diff.Unchanged:
			fmt.Fprintf(b, "%s%s: %s\n", indent, n.Key, stylishValue(n.Value, depth))
		case diff.Added:
			fmt.Fprintf(b, "%s+ %s: %s\n", markerIndent, n.Key, stylishValue(n.Value, depth))
		case diff.Removed:
			fmt.Fprintf(b, "%s- %s: %s\n", markerIndent, n.Key, stylishValue(n.Value, depth))
		case diff.Changed:
			fmt.Fprintf(b, "%s- %s: %s\n", markerIndent, n.Key, stylishValue(n.OldValue, depth))
			fmt.Fprintf(b, "%s+ %s: %s\n", markerIndent, n.Key, stylishValue(n.NewValue, depth))
		case diff.Nested:
			fmt.Fprintf(b, "%s%s: {\n", indent, n.Key)
			renderStylish(b, n.Children, depth+1)
			fmt.Fprintf(b, "%s}\n", indent)
		}
	}
}

func stylishValue(v any, depth int) string {
	switch val := v.(type) {
	case map[string]any:
		return stylishMap(val, depth)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", val)
	}
}

func stylishMap(m map[string]any, depth int) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	childIndent := strings.Repeat(" ", (depth+1)*stylishIndent)
	closeIndent := strings.Repeat(" ", depth*stylishIndent)

	var b strings.Builder
	b.WriteString("{\n")
	for _, k := range keys {
		fmt.Fprintf(&b, "%s%s: %s\n", childIndent, k, stylishValue(m[k], depth+1))
	}
	fmt.Fprintf(&b, "%s}", closeIndent)
	return b.String()
}
