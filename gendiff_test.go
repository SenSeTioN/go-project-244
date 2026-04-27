package code

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func loadExpected(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "fixture", name))
	if err != nil {
		t.Fatalf("read expected fixture %q: %v", name, err)
	}
	return strings.TrimRight(string(b), "\n")
}

func fixturePath(name string) string {
	return filepath.Join("testdata", "fixture", name)
}

func TestGenDiffNestedStylish(t *testing.T) {
	expected := loadExpected(t, "expected_nested_stylish.txt")

	cases := []struct {
		name   string
		file1  string
		file2  string
		format string
	}{
		{"json/json explicit stylish", "file1_nested.json", "file2_nested.json", "stylish"},
		{"yml/yml explicit stylish", "file1_nested.yml", "file2_nested.yml", "stylish"},
		{"json/yml mixed", "file1_nested.json", "file2_nested.yml", "stylish"},
		{"default format (empty string)", "file1_nested.json", "file2_nested.json", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GenDiff(fixturePath(tc.file1), fixturePath(tc.file2), tc.format)
			if err != nil {
				t.Fatalf("GenDiff returned error: %v", err)
			}
			if got != expected {
				t.Errorf("unexpected diff:\nwant:\n%s\n\ngot:\n%s", expected, got)
			}
		})
	}
}

func TestGenDiffNestedPlain(t *testing.T) {
	expected := loadExpected(t, "expected_nested_plain.txt")

	cases := []struct {
		name  string
		file1 string
		file2 string
	}{
		{"json/json", "file1_nested.json", "file2_nested.json"},
		{"yml/yml", "file1_nested.yml", "file2_nested.yml"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GenDiff(fixturePath(tc.file1), fixturePath(tc.file2), "plain")
			if err != nil {
				t.Fatalf("GenDiff returned error: %v", err)
			}
			if got != expected {
				t.Errorf("unexpected plain diff:\nwant:\n%s\n\ngot:\n%s", expected, got)
			}
		})
	}
}

func TestGenDiffNestedJSON(t *testing.T) {
	expected := loadExpected(t, "expected_nested_json.txt")

	cases := []struct {
		name  string
		file1 string
		file2 string
	}{
		{"json/json", "file1_nested.json", "file2_nested.json"},
		{"yml/yml", "file1_nested.yml", "file2_nested.yml"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GenDiff(fixturePath(tc.file1), fixturePath(tc.file2), "json")
			if err != nil {
				t.Fatalf("GenDiff returned error: %v", err)
			}

			var parsed map[string]any
			if err := json.Unmarshal([]byte(got), &parsed); err != nil {
				t.Fatalf("json output is not a valid JSON object: %v\n%s", err, got)
			}

			if tc.name == "json/json" && got != expected {
				t.Errorf("unexpected json diff:\nwant:\n%s\n\ngot:\n%s", expected, got)
			}
		})
	}
}

func TestGenDiffIdenticalFilesProducesNoMarkers(t *testing.T) {
	got, err := GenDiff(fixturePath("file1_nested.json"), fixturePath("file1_nested.json"), "stylish")
	if err != nil {
		t.Fatalf("GenDiff returned error: %v", err)
	}

	for _, marker := range []string{"  - ", "  + "} {
		if strings.Contains(got, marker) {
			t.Errorf("diff of identical files should not contain %q, got:\n%s", marker, got)
		}
	}
}

func TestGenDiffMissingFile(t *testing.T) {
	_, err := GenDiff(fixturePath("does_not_exist.json"), fixturePath("file2_nested.json"), "stylish")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestGenDiffUnknownFormat(t *testing.T) {
	_, err := GenDiff(fixturePath("file1_nested.json"), fixturePath("file2_nested.json"), "no-such-format")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}
