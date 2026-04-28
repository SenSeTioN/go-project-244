// Package parsers читает и разбирает конфигурационные файлы (JSON и YAML).
package parsers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Parse разбирает файл в map по его расширению (.json, .yml, .yaml).
func Parse(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %q: %w", path, err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return parseJSON(data)
	case ".yml", ".yaml":
		return parseYAML(data)
	default:
		return nil, fmt.Errorf("unsupported file extension %q", ext)
	}
}

func parseJSON(data []byte) (map[string]any, error) {
	result := make(map[string]any)
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parse json: %w", err)
	}
	return result, nil
}

func parseYAML(data []byte) (map[string]any, error) {
	result := make(map[string]any)
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}
	return result, nil
}
