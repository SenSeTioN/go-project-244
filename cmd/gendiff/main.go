// Команда gendiff сравнивает два конфигурационных файла и печатает дифф.
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"code"

	"github.com/urfave/cli/v2"
)

// stringifyFlag форматирует флаг для блока GLOBAL OPTIONS.
func stringifyFlag(f cli.Flag) string {
	var b strings.Builder
	switch fl := f.(type) {
	case *cli.StringFlag:
		b.WriteString("--")
		b.WriteString(fl.Name)
		b.WriteString(" string")
		for _, a := range fl.Aliases {
			b.WriteString(", -")
			b.WriteString(a)
			b.WriteString(" string")
		}
		fmt.Fprintf(&b, "\t%s (default: %q)", fl.Usage, fl.Value)
	case *cli.BoolFlag:
		b.WriteString("--")
		b.WriteString(fl.Name)
		for _, a := range fl.Aliases {
			b.WriteString(", -")
			b.WriteString(a)
		}
		fmt.Fprintf(&b, "\t%s", fl.Usage)
	}
	return b.String()
}

// newApp создаёт CLI-приложение с указанным writer'ом для вывода.
func newApp(out io.Writer) *cli.App {
	return &cli.App{
		Name:            "gendiff",
		Usage:           "Compares two configuration files and shows a difference.",
		UsageText:       "gendiff [global options] filepath1 filepath2",
		HideHelpCommand: true,
		Writer:          out,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "output format",
				Value:   "stylish",
			},
		},
		Action: runDiff,
	}
}

// runDiff — действие CLI: вызывает GenDiff и печатает результат.
func runDiff(c *cli.Context) error {
	if c.NArg() < 2 {
		return cli.ShowAppHelp(c)
	}

	result, err := code.GenDiff(c.Args().Get(0), c.Args().Get(1), c.String("format"))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(c.App.Writer, result)
	return err
}

func main() {
	cli.FlagStringer = stringifyFlag

	if err := newApp(os.Stdout).Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
