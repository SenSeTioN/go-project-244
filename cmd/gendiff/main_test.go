package main

import (
	"bytes"
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
