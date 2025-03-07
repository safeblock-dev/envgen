# envgen

`envgen` - это инструмент для генерации конфигурации окружения в различных форматах. Он позволяет описать конфигурацию в YAML файле и сгенерировать:
- Go структуры с тегами для [env](https://github.com/caarlos0/env)
- Файлы `.env.example`
- Документацию в формате Markdown
- Пользовательские форматы через собственные шаблоны

## Установка

```bash
go install github.com/safeblock-dev/envgen/cmd/envgen@latest
```

## Использование

```bash
# Базовая генерация
envgen -c config.yaml -o config.go -t template.tmpl

# Использование шаблона по URL
envgen -c config.yaml -o config.go -t https://raw.githubusercontent.com/user/repo/template.tmpl

# Просмотр версии
envgen version

# Справка по использованию
envgen --help
```

## Встроенные шаблоны

### go-env
Генерирует Go код с структурами для работы с переменными окружения:
- Поддержка вложенных структур
- Теги для [env](https://github.com/caarlos0/env)
- Валидация обязательных полей
- Поддержка пользовательских типов
- Документация в формате godoc

Примеры использования можно найти в тестах в директории `templates_tests/go-env`.

### example
Создает файл `.env.example` с примерами всех переменных:
- Значения по умолчанию
- Комментарии с описанием
- Группировка по разделам
- Поддержка префиксов

Примеры использования можно найти в тестах в директории `templates_tests/example`.

### markdown
Генерирует документацию в формате Markdown:
- Полное описание всех переменных
- Таблицы с типами и значениями по умолчанию
- Примеры использования
- Инструкции по установке
- Доступен на русском (markdown-ru) и английском языках

Примеры документации будут доступны в следующих релизах.

Исходные шаблоны:
- [templates/go-env](templates/go-env) - шаблон для Go кода
- [templates/example](templates/example) - шаблон для .env.example
- [templates/markdown](templates/markdown) - шаблон для английской документации
- [templates/markdown-ru](templates/markdown-ru) - шаблон для русской документации

## Шаблоны

Шаблоны можно указывать двумя способами:
1. Путь к локальному файлу: `-t ./templates/config.tmpl`
2. URL: `-t https://raw.githubusercontent.com/user/repo/template.tmpl`

В шаблонах доступны следующие функции:
- Преобразование строк: `title`, `upper`, `lower`, `camel`, `snake`, `kebab`, `pascal`
- Работа с типами: `toString`, `toInt`, `toBool`
- Дата и время: `now`, `formatTime`, `date`, `datetime`
- Условные операторы: `default`, `coalesce`, `ternary`
- Операции со строками: `contains`, `hasPrefix`, `hasSuffix`, `replace`, `trim`, `join`, `split`
- Работа с путями: `getDirName`, `getFileName`, `getFileExt`, `joinPaths`

Подробные примеры использования этих функций можно найти в тестах в директории `pkg/envgen/templatefuncs`.

## Примеры

- [Базовый пример](examples/basic) - стандартное использование с несколькими группами
- [Минимальный пример](examples/minimal) - простейшая конфигурация
- [Пользовательский шаблон](examples/custom) - пример с собственным шаблоном
- [Go Generate](examples/gogen) - использование с `go:generate`

## Конфигурация

Конфигурация описывается в YAML файле и поддерживает:
- Группировку переменных
- Префиксы для групп
- Описание и примеры для каждой переменной
- Значения по умолчанию
- Пометки обязательных полей
- Дополнительные опции для кастомизации

Примеры конфигурации можно найти в тестах в директории `templates_tests`.

## Лицензия

MIT 