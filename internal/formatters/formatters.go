// Package formatters рендерит дерево диффа в строку (stylish, plain, json).
package formatters

import (
	"fmt"

	"code/internal/diff"
)

// DefaultFormat — форматер по умолчанию.
const DefaultFormat = "stylish"

// Format рендерит дерево в строку. Пустое name означает DefaultFormat.
func Format(name string, nodes []diff.Node) (string, error) {
	if name == "" {
		name = DefaultFormat
	}
	switch name {
	case "stylish":
		return formatStylish(nodes), nil
	case "plain":
		return formatPlain(nodes), nil
	case "json":
		return formatJSON(nodes)
	default:
		return "", fmt.Errorf("unknown format %q", name)
	}
}
