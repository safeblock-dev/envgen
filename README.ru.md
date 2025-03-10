# envgen

**`envgen` – мощный инструмент для генерации конфигурации окружения из YAML**  

Создаёт типобезопасные Go-структуры, `.env.example`, документацию и любые другие файлы по пользовательским шаблонам

### Преимущества

- 🔒 **Типобезопасность**: Валидация типов на этапе компиляции
- 🔄 **Автоматическая генерация**: Документация и примеры всегда синхронизированы с кодом
- 🎨 **Любой шаблон**: Поддержка пользовательских шаблонов и форматов
- 🛠 **Гибкая настройка**: Используйте простые настройки для кастомизации
- 📝 **Автодокументация**: Markdown-документация генерируется автоматически
- 🔍 **Прозрачность**: Понятная структура конфигурации в YAML-формате

## Возможности

- Множество выходных форматов:
  - Go-структуры с тегами `env`
  - Пример файлов окружения (`.env.example`)
  - Документация в формате Markdown
- Настраиваемые шаблоны
- Возможность использовать собственные шаблоны

## Установка

```bash
go install github.com/safeblock-dev/envgen/cmd/envgen@latest
```

## Быстрый старт

1. Создайте файл конфигурации `config.yaml`:
```yaml
# Пример конфигурации для генерации Go-пакета
options:
  go_package: config  # Имя пакета для сгенерированного кода

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
        example: "0.0.0.0"
      
      - name: env
        type: Environment
        description: Окружение
        default: "development"
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

### Опции

Опции позволяют настраивать и модифицировать информацию в шаблоне. Для разных шаблонов используются разные опции.

```yaml
options:
  go_package: mypkg
```

В данном примере показана опция, которая установит имя пакета для `go-env` шаблона. Вы также можете задать опции в группе (`group`) и в отдельном поле (`field`).

### Группы

Группы организуют связанные поля конфигурации:

```yaml
groups:
  - name: Database     # Обязательное: имя группы
    description: Настройки базы данных # Опциональное: описание группы
    prefix: DB_         # Опциональное: префикс для переменных окружения
    options:            # Опциональное: параметры группы
      go_name: DBConfig # Опциональное: любая опция для шаблона
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
  - name: url                   # Обязательное: имя переменной окружения
    type: string                # Обязательное: тип поля (встроенный или пользовательский)
    description: API endpoint   # Опциональное: описание поля
    default: "http://127.0.0.1" # Опциональное: значение по умолчанию
    required: true              # Опциональное: является ли поле обязательным
    example: "http://test.com"  # Опциональное: пример значения для документации
    options:                    # Опциональное: дополнительные параметры поля
      go_name: "URL"            # Опциональное: любая опция для шаблона
```

### Типы

Типы позволяют определять пользовательские типы, добавлять контекст к типу и переиспользовать их:

```yaml
types:
  - name: Duration        # Обязательное: имя типа для ссылок в полях
    type: time.Duration   # Обязательное: определение типа (встроенный или пользовательский)
    import: time          # Опциональное: путь импорта для пользовательских типов
    description: Интервал # Опциональное: описание типа
    values:               # Опциональное: возможные значения для документации
      - 1s
      - 1m
```

Вы можете создать несколько похожих типов:

```yaml
types:
  - name: AppENV
    type: string
    description: Имя окружения
    values: ["prod", "dev"]
  - name: MediaURL
    type: string
    description: Ссылка на media источник
```

Чтобы использовать созданные типы, нужно указать их имя (`name`) в качестве значения типа (`type`) в описании поля (`field`):

```yaml
fields:
  - name: github                  
    type: AppURL                # Указываем имя type
    example: "http://github.com/safeblock-dev" 
  - name: twitter                  
    type: AppURL                # Тип можно использовать несколько раз
    example: "http://x.com/safeblock" 
```

## Продвинутые возможности

### Композитные конфигурации

Вы можете игнорировать типы и группы при генерации:

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

Это особенно полезно, когда у вас есть структуры, которые вы не хотите показывать, например, в `.env.example`.

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

### Опции Go-шаблона

Шаблон `go-env` поддерживает глобальные опции:

```yaml
options:
  go_package: config # Обязательное поле
  go_generate: |
    # Генерация конфигурации
    //go:generate envgen -c {{ ConfigPath }} -o {{ OutputPath }} -t {{ TemplatePath }}
    # Генерация документации
    //go:generate envgen -c {{ ConfigPath }} -o docs/{{ OutputPath }} -t markdown
  go_meta: |
    // Версия: v0.1.2
    // Шаблон: {{ TemplatePath }}
```

Опция `go_package` является обязательной для шаблона `go-env`. Если значение не указано, `envgen` попытается использовать имя папки из флага `out`, но это считается плохой практикой, поскольку если путь будет вида `config.go`, то имя пакета будет установлено как `.`, что приведет к ошибке компиляции.

Опция `go_generate` позволяет указать пользовательские команды для генерации кода. Если эта опция не указана, используется команда по умолчанию.

Опция `go_meta` добавит дополнительную информацию после `go_generate` блока.

В опциях (`go_generate`, `go_meta`) доступны следующие специальные ключи:
- `{{ ConfigPath }}` - выведет путь к файлу конфигурации
- `{{ OutputPath }}` - выведет путь к выходному файлу
- `{{ TemplatePath }}` - выведет путь к шаблону

Также доступны специальные опции для групп (`groups`) и полей (`fields`):

```yaml
groups:
  - name: App
    description: Application settings
    options:
      go_name: CustomAppConfig
    fields:
      - name: debug_mode
        type: bool
        description: Enable debug mode
        options:
          go_name: IsDebug
```

Результат выполнения:

```go
// Имя App изменено на CustomAppConfig
type CustomAppConfig struct {
   // Имя DebugMode (debug_mode в файле конфигурации) изменено на IsDebug
	IsDebug bool `env:"DEBUG_MODE" envDefault:"false"`
}
```

### Опции Markdown шаблона

Шаблон `markdown` поддерживает глобальные опции:

```yaml
options:
  md_title: Заголовок Markdown файла
  md_description: |
    Любое дополнительное описание
```

Эти опции являются дополнительными и не являются обязательными.

### Опции Example шаблона

Шаблон не использует каких-либо специальных опций.

## Разработка

### Структура проекта

```
.
├── cmd/envgen/             # CLI приложение
├── pkg/envgen/             # Основной пакет
│   ├── config.go           # Типы конфигурации и валидация
│   ├── envgen.go           # Основная логика генерации
│   ├── template.go         # Загрузка шаблонов
│   ├── template_context.go # Контекст шаблона и функции
│   └── templatefuncs/      # Вспомогательные функции шаблонов
├── templates/              # Встроенные шаблоны
└── templates_tests/        # Тесты и примеры
```

### Запуск тестов

```bash
go test ./...
```

Обновление golden-файлов для тестов шаблонов:
```bash
UPDATE_GOLDEN=1 go test ./templates_tests
```

## Часто задаваемые вопросы

### Как добавить свой шаблон?

Создайте файл шаблона с расширением `.tmpl` или `.tpl`. Используйте синтаксис Go templates и доступные функции из контекста. Пример простого шаблона:

```go
// Файл: custom.tmpl
{{- range $group := .Groups }}
# {{ $group.Description }}
{{- range $field := $group.Fields }}
{{ $field.Name | upper }}={{ $field.Default }}  # {{ $field.Description }}
{{- end }}

{{- end }}
```

Сгенерируйте шаблон:
```bash
envgen -c config.yaml -o custom.txt -t ./custom.tmpl
```

Результат будет выглядеть так:
```ini
# Настройки веб-сервера
PORT=8080  # Порт сервера
HOST=localhost  # Хост сервера
ENV=development  # Окружение

# Настройки базы данных
DB_HOST=localhost  # Хост базы данных
DB_PORT=5432  # Порт базы данных
```

### Как использовать пользовательские типы?

1. Определите тип в конфигурации:
```yaml
types:
  - name: CustomType
    type: your/pkg.Type
    import: your/pkg
```

2. Используйте его в полях:
```yaml
fields:
  - name: custom_field
    type: CustomType
```

### Какие функции доступны в шаблонах?

В шаблонах доступны следующие встроенные функции:

- Функции для работы со строками:
  - `title` - преобразование первой буквы в верхний регистр
  - `upper` - преобразование в верхний регистр
  - `lower` - преобразование в нижний регистр
  - `pascal` - преобразование в PascalCase
  - `camel` - преобразование в camelCase
  - `snake` - преобразование в snake_case
  - `kebab` - преобразование в kebab-case
  - `append` - добавление строки в конец
  - `uniq` - удаление дубликатов
  - `slice` - получение подстроки
  - `contains` - проверка наличия подстроки
  - `hasPrefix` - проверка префикса
  - `hasSuffix` - проверка суффикса
  - `replace` - замена подстроки
  - `trim` - удаление пробельных символов
  - `join` - объединение строк
  - `split` - разделение строки

- Функции для работы с типами:
  - `toString` - преобразование в строку
  - `toInt` - преобразование в целое число
  - `toBool` - преобразование в логическое значение
  - `findType` - поиск информации о типе
  - `getImports` - получение списка импортов
  - `typeImport` - получение импорта для типа

- Функции для работы с датой и временем:
  - `now` - текущее время
  - `formatTime` - форматирование времени
  - `date` - текущая дата (ГГГГ-ММ-ДД)
  - `datetime` - текущие дата и время (ГГГГ-ММ-ДД ЧЧ:ММ:СС)

- Условные операции:
  - `default` - значение по умолчанию
  - `coalesce` - первое непустое значение
  - `ternary` - тернарный оператор
  - `hasOption` - проверка наличия опции
  - `hasGroupOption` - проверка наличия опции в группе
  - `getOption` - получение значения опции
  - `getGroupOption` - получение значения опции из группы

- Функции для работы с путями:
  - `getDirName` - получение имени директории
  - `getFileName` - получение имени файла
  - `getFileExt` - получение расширения файла
  - `joinPaths` - объединение путей
  - `getConfigPath` - путь к файлу конфигурации
  - `getOutputPath` - путь к выходному файлу
  - `getTemplatePath` - путь к файлу шаблона

Пример использования:
```go
{{ $name := "my_variable" }}
{{ $name | pascal }}  // Результат: MyVariable
{{ $name | upper }}   // Результат: MY_VARIABLE

{{ if hasOption "go_package" }}
package {{ getOption "go_package" }}
{{ end }}

// Работа с датой и временем
{{ datetime }}  // Результат: 2024-03-21 15:04:05

// Условные операции
{{ $value := "test" | default "default_value" }}

// Работа с путями
{{ $dir := getFileName "/path/to/file.txt" }}  // Результат: file.txt
```

## Лицензия

MIT