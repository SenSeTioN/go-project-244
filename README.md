### Hexlet tests and linter status:

[![Actions Status](https://github.com/SenSeTioN/go-project-244/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/SenSeTioN/go-project-244/actions)
[![CI](https://github.com/SenSeTioN/go-project-244/actions/workflows/ci.yml/badge.svg)](https://github.com/SenSeTioN/go-project-244/actions/workflows/ci.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=SenSeTioN_go-project-244&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=SenSeTioN_go-project-244)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=SenSeTioN_go-project-244&metric=coverage)](https://sonarcloud.io/summary/new_code?id=SenSeTioN_go-project-244)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=SenSeTioN_go-project-244&metric=bugs)](https://sonarcloud.io/summary/new_code?id=SenSeTioN_go-project-244)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=SenSeTioN_go-project-244&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=SenSeTioN_go-project-244)

## О проекте

`gendiff` — CLI-утилита для сравнения двух конфигурационных файлов.
Поддерживает форматы входа `JSON`, `YML`/`YAML` и три формата вывода:
`stylish`, `plain`, `json`. Файлы сравниваются как структуры данных, а
не как текст, поэтому порядок ключей не влияет на результат.

## Команды утилиты

| Команда                              | Описание                                |
| ------------------------------------ | --------------------------------------- |
| `./bin/gendiff -h`                   | Показать справку                        |
| `./bin/gendiff file1 file2`          | Дифф в формате `stylish` (по умолчанию) |
| `./bin/gendiff -f plain file1 file2` | Дифф в формате `plain`                  |
| `./bin/gendiff -f json file1 file2`  | Дифф в формате `json`                   |

## Команды Makefile

| Команда              | Описание                                         |
| -------------------- | ------------------------------------------------ |
| `make build`         | Сборка бинарника в `bin/gendiff`                 |
| `make run`           | Сборка и запуск без аргументов (покажет справку) |
| `make test`          | Запуск тестов                                    |
| `make test-coverage` | Тесты с покрытием → `coverage.out`               |
| `make lint`          | Запуск `golangci-lint`                           |
