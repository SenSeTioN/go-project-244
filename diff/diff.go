// Package diff строит внутреннее дерево различий между двумя структурами.
package diff

import (
	"reflect"
	"sort"
)

// Status описывает изменение ключа в дереве диффа.
type Status string

// Возможные значения Status.
const (
	Unchanged Status = "unchanged" // значения совпадают
	Added     Status = "added"     // ключ только во втором файле
	Removed   Status = "removed"   // ключ только в первом файле
	Changed   Status = "changed"   // значения различаются
	Nested    Status = "nested"    // оба значения — объекты, сравнение вглубь
)

// Node — узел дерева диффа для одного ключа.
type Node struct {
	Key      string
	Status   Status
	Value    any
	OldValue any
	NewValue any
	Children []Node
}

// Build сравнивает две map и возвращает отсортированное дерево диффа.
func Build(data1, data2 map[string]any) []Node {
	keys := mergedKeys(data1, data2)
	nodes := make([]Node, 0, len(keys))
	for _, k := range keys {
		v1, ok1 := data1[k]
		v2, ok2 := data2[k]
		switch {
		case ok1 && !ok2:
			nodes = append(nodes, Node{Key: k, Status: Removed, Value: v1})
		case !ok1 && ok2:
			nodes = append(nodes, Node{Key: k, Status: Added, Value: v2})
		default:
			m1, isMap1 := v1.(map[string]any)
			m2, isMap2 := v2.(map[string]any)
			switch {
			case isMap1 && isMap2:
				nodes = append(nodes, Node{Key: k, Status: Nested, Children: Build(m1, m2)})
			case reflect.DeepEqual(v1, v2):
				nodes = append(nodes, Node{Key: k, Status: Unchanged, Value: v1})
			default:
				nodes = append(nodes, Node{Key: k, Status: Changed, OldValue: v1, NewValue: v2})
			}
		}
	}
	return nodes
}

func mergedKeys(data1, data2 map[string]any) []string {
	seen := make(map[string]struct{}, len(data1)+len(data2))
	for k := range data1 {
		seen[k] = struct{}{}
	}
	for k := range data2 {
		seen[k] = struct{}{}
	}

	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
