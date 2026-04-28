// Package code предоставляет функцию GenDiff — публичный API библиотеки.
package code

import (
	"code/internal/diff"
	"code/internal/formatters"
	"code/internal/parsers"
)

// GenDiff возвращает дифф двух файлов в указанном формате.
// Пустой format означает форматер по умолчанию.
func GenDiff(filepath1, filepath2, format string) (string, error) {
	data1, err := parsers.Parse(filepath1)
	if err != nil {
		return "", err
	}

	data2, err := parsers.Parse(filepath2)
	if err != nil {
		return "", err
	}

	tree := diff.Build(data1, data2)
	return formatters.Format(format, tree)
}
