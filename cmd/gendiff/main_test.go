package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/urfave/cli/v2"
)

func fixturePath(t *testing.T, name string) string {
	t.Helper()
	p, err := filepath.Abs(filepath.Join("..", "..", "testdata", "fixture", name))
	if err != nil {
		t.Fatalf("abs path: %v", err)
	}
	return p
}

func loadExpected(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(fixturePath(t, name))
	if err != nil {
		t.Fatalf("read expected fixture %q: %v", name, err)
	}
	return strings.TrimRight(string(b), "\n")
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
			got, err := GenDiff(fixturePath(t, tc.file1), fixturePath(t, tc.file2), tc.format)
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
			got, err := GenDiff(fixturePath(t, tc.file1), fixturePath(t, tc.file2), "plain")
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
			got, err := GenDiff(fixturePath(t, tc.file1), fixturePath(t, tc.file2), "json")
			if err != nil {
				t.Fatalf("GenDiff returned error: %v", err)
			}

			var parsed []map[string]any
			if err := json.Unmarshal([]byte(got), &parsed); err != nil {
				t.Fatalf("json output is not valid JSON: %v\n%s", err, got)
			}

			if tc.name == "json/json" && got != expected {
				t.Errorf("unexpected json diff:\nwant:\n%s\n\ngot:\n%s", expected, got)
			}
		})
	}
}

func TestGenDiffIdenticalFilesProducesNoMarkers(t *testing.T) {
	file1 := fixturePath(t, "file1_nested.json")

	got, err := GenDiff(file1, file1, "stylish")
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
	_, err := GenDiff(fixturePath(t, "does_not_exist.json"), fixturePath(t, "file2_nested.json"), "stylish")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestGenDiffUnknownFormat(t *testing.T) {
	_, err := GenDiff(fixturePath(t, "file1_nested.json"), fixturePath(t, "file2_nested.json"), "no-such-format")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func runApp(t *testing.T, args ...string) string {
	t.Helper()
	var buf bytes.Buffer
	app := newApp(&buf)
	if err := app.Run(append([]string{"gendiff"}, args...)); err != nil {
		t.Fatalf("app run: %v", err)
	}
	return buf.String()
}

func TestRunDiffStylish(t *testing.T) {
	out := runApp(t,
		fixturePath(t, "file1_nested.json"),
		fixturePath(t, "file2_nested.json"),
	)
	if !strings.Contains(out, "+ follow: false") {
		t.Errorf("expected stylish diff to contain added follow, got:\n%s", out)
	}
	if !strings.Contains(out, "- setting2: 200") {
		t.Errorf("expected stylish diff to contain removed setting2, got:\n%s", out)
	}
}

func TestRunDiffFormatPlain(t *testing.T) {
	out := runApp(t,
		"--format", "plain",
		fixturePath(t, "file1_nested.json"),
		fixturePath(t, "file2_nested.json"),
	)
	if !strings.Contains(out, "Property 'common.follow' was added with value: false") {
		t.Errorf("expected plain output with follow line, got:\n%s", out)
	}
}

func TestRunDiffMissingArgumentsShowsHelp(t *testing.T) {
	out := runApp(t, fixturePath(t, "file1_nested.json"))
	if !strings.Contains(out, "NAME:") || !strings.Contains(out, "gendiff") {
		t.Errorf("expected help output when args are missing, got:\n%s", out)
	}
}

func TestRunDiffInvalidFormatReturnsError(t *testing.T) {
	var buf bytes.Buffer
	app := newApp(&buf)
	err := app.Run([]string{
		"gendiff",
		"--format", "no-such",
		fixturePath(t, "file1_nested.json"),
		fixturePath(t, "file2_nested.json"),
	})
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestStringifyFlagStringFlag(t *testing.T) {
	s := stringifyFlag(&cli.StringFlag{
		Name:    "format",
		Aliases: []string{"f"},
		Usage:   "output format",
		Value:   "stylish",
	})
	for _, want := range []string{"--format string", "-f string", "output format", `"stylish"`} {
		if !strings.Contains(s, want) {
			t.Errorf("expected %q in %q", want, s)
		}
	}
}

func TestStringifyFlagBoolFlag(t *testing.T) {
	s := stringifyFlag(&cli.BoolFlag{Name: "help", Aliases: []string{"h"}, Usage: "show help"})
	for _, want := range []string{"--help", "-h", "show help"} {
		if !strings.Contains(s, want) {
			t.Errorf("expected %q in %q", want, s)
		}
	}
}
