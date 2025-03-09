# envgen

`envgen` - это гибкий инструмент для генерации кода конфигурации окружения. Он создает типобезопасные структуры конфигурации и документацию на основе YAML-определений.

## Возможности

- Множество выходных форматов:
  - Go-структуры с тегами `env`
  - Шаблоны файлов окружения
  - Документация в формате Markdown
- Система типов с поддержкой пользовательских определений
- Организация конфигурации на основе групп
- Поддержка префиксов для переменных окружения
- Настраиваемые шаблоны
- Игнорирование типов и групп для композитных конфигураций
- Расширенная валидация и отчетность об ошибках

## Установка

```bash
go install github.com/safeblock-dev/envgen/cmd/envgen@latest
```

## Быстрый старт

1. Создайте файл конфигурации `config.yaml`:
```yaml
# Пример конфигурации для генерации Go-пакета
options:
  go_package: myapp/config  # Имя пакета для сгенерированного кода

types:
  - name: Environment
    type: string
    description: Окружение приложения (development, staging, production)
    values:
      - development
      - staging
      - production

groups:
  - name: Server
    description: Настройки веб-сервера
    prefix: SERVER_  # Префикс для переменных окружения
    fields:
      - name: port
        type: int
        description: Порт сервера
        default: "8080"
        required: true
        example: "9000"
      
      - name: host
        type: string
        description: Хост сервера
        default: "localhost"
        required: true
        example: "0.0.0.0"
      
      - name: env
        type: Environment
        description: Окружение
        default: "development"
        required: true
        example: "production"
```

2. Сгенерируйте Go-код:
```bash
envgen -c config.yaml -o config.go -t go-env
```

### Флаги командной строки

Инструмент поддерживает следующие флаги:

- `-c, --config`: Путь к входному YAML-файлу конфигурации (обязательный)
- `-o, --out`: Путь к выходному файлу (обязательный)
- `-t, --template`: Путь к файлу шаблона или URL (обязательный)
- `--ignore-types`: Список типов для игнорирования через запятую
- `--ignore-groups`: Список групп для игнорирования через запятую

Примеры:
```bash
# Генерация с использованием локального шаблона
envgen -c config.yaml -o config.go -t ./templates/config.tmpl

# Генерация с использованием шаблона из URL
envgen --config config.yaml --out config.go --template https://raw.githubusercontent.com/user/repo/template.tmpl

# Генерация с игнорированием определенных типов и групп
envgen -c config.yaml -o config.go -t ./templates/config.tmpl --ignore-types Duration,URL --ignore-groups Database

# Показать версию
envgen version
```

Это создаст файл `config.go` с типобезопасными структурами для работы с конфигурацией. Сгенерированный код будет использовать переменные окружения с префиксом `SERVER_` (например, `SERVER_PORT`, `SERVER_HOST`, `SERVER_ENV`).

## Формат конфигурации

### Типы

Типы позволяют определять пользовательские типы с валидацией и документацией:

```yaml
types:
  - name: Duration      # Обязательное: имя типа для ссылок в полях
    type: time.Duration # Обязательное: определение типа (встроенный или пользовательский)
    import: time       # Опциональное: путь импорта для пользовательских типов
    description: Временной интервал # Опциональное: описание типа
    values:           # Опциональное: возможные значения для документации
      - 1s
      - 1m
```

### Группы

Группы организуют связанные поля конфигурации:

```yaml
groups:
  - name: Database     # Обязательное: имя группы
    description: Настройки базы данных # Опциональное: описание группы
    prefix: DB_         # Опциональное: префикс для переменных окружения
    options:            # Опциональное: параметры группы
      go_name: DBConfig # Опциональное: переопределение имени структуры
    fields:             # Обязательное: должно быть определено хотя бы одно поле
      - name: host
        type: string
        description: Хост базы данных
        required: true
        default: localhost
```

### Поля

Поля представляют отдельные переменные окружения:

```yaml
fields:
  - name: port        # Обязательное: имя переменной окружения
    type: int        # Обязательное: тип поля (встроенный или пользовательский)
    description: Порт # Опциональное: описание поля
    default: "8080"  # Опциональное: значение по умолчанию
    required: true   # Опциональное: является ли поле обязательным
    example: "9000"  # Опциональное: пример значения для документации
    options:         # Опциональное: дополнительные параметры поля
      import: "custom/pkg" # Опциональное: путь импорта для пользовательских типов
      name_field: Port    # Опциональное: переопределение имени поля структуры
```

## Продвинутые возможности

### Композитные конфигурации

Вы можете использовать группы как типы и игнорировать их при генерации:

```yaml
groups:
  - name: Postgres
    description: Настройки PostgreSQL
    prefix: PG_
    fields:
      - name: host
        type: string
        default: localhost

  - name: Redis
    description: Настройки Redis
    prefix: REDIS_
    fields:
      - name: port
        type: int
        default: "6379"

  - name: Webserver
    description: Конфигурация веб-сервера
    fields:
      - name: db
        type: Postgres
        options:
          name_field: DB
      - name: cache
        type: Redis
```

Генерация только конфигураций баз данных:
```bash
envgen -c config.yaml -o config.go -t go-env --ignore-groups Webserver
```

### Шаблоны

Инструмент включает три встроенных шаблона:

- `go-env`: Генерирует Go-структуры с тегами `env`
- `example`: Создает шаблоны файлов `.env`
- `markdown`: Создает документацию в формате Markdown

Примеры использования встроенных шаблонов:

```bash
# Генерация Go-структур
envgen -c config.yaml -o config.go -t go-env

# Генерация шаблона .env файла
envgen -c config.yaml -o .env.example -t example

# Генерация документации
envgen -c config.yaml -o config.md -t markdown
```

### Использование сгенерированного кода

После генерации кода вы можете использовать его в вашем Go-приложении:

```go
package main

import (
    "fmt"
    "log"
    "os"

    "myapp/config"
)

func main() {
    var cfg config.ServerConfig
    if err := env.Parse(&cfg); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Сервер запустится на %s:%d\n", cfg.Host, cfg.Port)
    fmt.Printf("Окружение: %s\n", cfg.Env)
}
```

Сгенерированный код использует пакет `env` для парсинга переменных окружения. Убедитесь, что он добавлен в ваши зависимости:

```bash
go get github.com/caarlos0/env/v11
```

### Опции Go-шаблона

Шаблон `go-env` поддерживает дополнительные опции для полей, которые управляют обработкой переменных окружения:

```yaml
fields:
  - name: api_key
    type: string
    description: API ключ из файла
    required: true      # Глобальная опция: добавляет тег ,required
    options:
      go_file: true     # Чтение значения из файла
      go_expand: true   # Включить подстановку переменных окружения
      go_init: true     # Инициализировать nil указатели
      go_notEmpty: true # Ошибка если значение пустое
      go_unset: true    # Удалить переменную окружения после чтения
```

#### Доступные Go-опции

- `go_file`: Указывает, что значение должно быть прочитано из файла, указанного в переменной окружения
- `go_expand`: Включает подстановку переменных окружения в значениях (например, `FOO_${BAR}`)
- `go_init`: Инициализирует указатели, которые иначе были бы nil
- `go_notEmpty`: Возвращает ошибку, если переменная окружения пуста
- `go_unset`: Удаляет переменные окружения после их прочтения

Эти опции реализованы с использованием тегов пакета `github.com/caarlos0/env/v11`.

## Разработка

### Структура проекта

```