package formatters

import (
	"encoding/json"
	"fmt"

	"code/diff"
)

func formatJSON(nodes []diff.Node) (string, error) {
	data, err := json.MarshalIndent(toJSONNodes(nodes), "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal json diff: %w", err)
	}
	return string(data), nil
}

func toJSONNodes(nodes []diff.Node) []map[string]any {
	out := make([]map[string]any, 0, len(nodes))
	for _, n := range nodes {
		entry := map[string]any{
			"key":    n.Key,
			"status": string(n.Status),
		}
		switch n.Status {
		case diff.Nested:
			entry["children"] = toJSONNodes(n.Children)
		case diff.Changed:
			entry["oldValue"] = n.OldValue
			entry["newValue"] = n.NewValue
		case diff.Added, diff.Removed, diff.Unchanged:
			entry["value"] = n.Value
		}
		out = append(out, entry)
	}
	return out
}
