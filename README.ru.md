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
      - name: Port
        type: int
        description: Порт сервера
        default: "8080"
        required: true
        example: "9000"
      
      - name: Host
        type: string
        description: Хост сервера
        default: "localhost"
        example: "0.0.0.0"
      
      - name: ENV
        type: Environment
        description: Окружение
        default: "development"
        example: "production"
```

2. Сгенерируйте Go-код:
```bash
envgen gen -c config.yaml -o config.go -t go-env
```

### Команды

`envgen` поддерживает следующие команды:

- `gen` (или `generate`): Генерация файлов конфигурации
  - `-c, --config`: Путь к входному YAML-файлу конфигурации (обязательный)
  - `-o, --out`: Путь к выходному файлу (обязательный)
  - `-t, --template`: Путь к файлу шаблона или URL (обязательный)
  - `--ignore-types`: Список типов для игнорирования через запятую
  - `--ignore-groups`: Список групп для игнорирования через запятую

- `ls` (или `templates`, `list`): Показать список доступных стандартных шаблонов

- `version`: Показать версию программы

Примеры:
```bash
# Генерация с использованием локального шаблона
envgen gen -c config.yaml -o config.go -t ./templates/config.tmpl

# Генерация с использованием шаблона из URL
envgen gen --config config.yaml --out config.go --template https://raw.githubusercontent.com/user/repo/template.tmpl

# Генерация с игнорированием определенных типов и групп
envgen gen -c config.yaml -o config.go -t ./templates/config.tmpl --ignore-types Duration,URL --ignore-groups Database

# Показать список доступных шаблонов
envgen ls

# Показать версию
envgen version

# Генерация Go-структур
envgen gen -c config.yaml -o config.go -t go-env

# Генерация шаблона .env файла
envgen gen -c config.yaml -o .env.example -t example

# Генерация документации
envgen gen -c config.yaml -o config.md -t markdown
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
      - name: Host
        type: string
        description: Хост базы данных
        required: true
        default: localhost
```

### Поля

Поля представляют отдельные переменные окружения:

```yaml
fields:
  - name: URL                   # Обязательное: имя переменной окружения
    type: string                # Обязательное: тип поля (встроенный или пользовательский)
    description: API endpoint   # Опциональное: описание поля
    default: "http://127.0.0.1" # Опциональное: значение по умолчанию
    required: true              # Опциональное: является ли поле обязательным
    example: "http://test.com"  # Опциональное: пример значения для документации
    options:                    # Опциональное: дополнительные параметры поля
      go_name: "GitURL"         # Опциональное: любая опция для шаблона
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
  - name: Github                  
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
      - name: Host
        type: string
        default: localhost

  - name: Redis
    description: Настройки Redis
    prefix: REDIS_
    fields:
      - name: Port
        type: int
        default: "6379"

  - name: Webserver
    description: Конфигурация веб-сервера
    fields:
      - name: DB
        type: Postgres
      - name: Cache
        type: Redis
```

Генерация только конфигураций баз данных:
```bash
envgen gen -c config.yaml -o config.go -t go-env --ignore-groups Webserver
```

Это особенно полезно, когда у вас есть структуры, которые вы не хотите показывать, например, в `.env.example`.

### Шаблоны

Инструмент включает три встроенных шаблона:

- `go-env`: Генерирует Go-структуры с тегами `env`
- `go-env-example`: Создает `.env.example` шаблоны с учетом `go-env` тегов (опций)
- `example`: Создает `.env.example` шаблоны
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
  go_package: config # Не обязательное поле
  go_meta: |
    # Генерация конфигурации
    {{ goCommentGenerate "" "" "" }}
    # Генерация документации
    {{ goCommentGenerate "" "docs/README.md" "../../templates/markdown" }}
    // Версия: v0.1.2
```

Если значение `go_package` не указано, `envgen` попытается использовать имя папки из флага `out`.

Опция `go_meta` позволяет указать пользовательские команды для генерации кода. Если эта опция не указана, используется команда по умолчанию. Если вы не хотите, чтобы был вывод `//go:generate`, оставьте поле `go_meta` пустым.

Опция `go_meta` позволяет вызывать любые шаблонные функции из файла [funcs.go](pkg/envgen/funcs.go) (например `title`, `upper` и т.п.).

Также доступны специальные опции по настройки имени для групп (`groups`) и полей (`fields`):

```yaml
groups:
  - name: App
    description: Application settings
    options:
      go_name: CustomAppConfig
    fields:
      - name: DebugMode
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

Дополнительные опции для настройки тегов `env` для групп и полей:

- `go_skip_env_tag` - отключает генерацию тега `env`

Пример использования:

```yaml
  - name: NoEnvTag
    description: Группа без тегов env
    options:
      go_skip_env_tag: true
    fields:
      - name: Sentry
        type: SentryConfig
      - name: GRPC_Port
        type: int
        default: "8002"
      - name: HTTP_Port
        type: int
        default: "8001"
        options:
          go_name: HttpPort
          go_tags: env:"NOT_SKIPPED"
    
  - name: CustomEnvTags
    description: Выборочное применение тегов env
    fields:
      - name: NotSkipped
        type: string
      - name: debug
        type: bool
        options:
          go_skip_env_tag: true
      - name: Port
        type: int
        options:
          go_skip_env_tag: true
          go_env_options: will_be_ignored
          go_tags: env:"NOT_SKIPPED,required,notEmpty"
```

опции только для полей:

- `go_include` - если true, использует встраивание структур Go (struct embedding)
- `go_env_options` - позволяет добавить дополнительные опции в тег `env`. Например: `file`, `unset`, `notEmpty` и другие опции. Все опции передаются напрямую в теги без дополнительной валидации.
- `go_tags` - позволяет добавить дополнительные теги для структуры. Поддерживает указание любых тегов без ограничений. При использовании с пакетом [`env`](github.com/caarlos0/env/v11) часто используются:
  - `envSeparator` - разделитель для слайсов
  - `envKeyValSeparator` - разделитель для ключей и значений в мапах

Пример использования:

```yaml
groups:
  - name: Webserver
    fields:
      - name: Config
        type: Config
        options:
          go_skip_env_tag: true
          go_include: true # Встраивание структуры

      - name: ApiKey
        type: string
        description: API ключ, который будет очищен после чтения
        required: true
        options:
          go_env_options: unset,notEmpty  # Очистить после чтения и проверить на пустоту
        example: "secret-key"

      - name: Tags
        type: "[]string"
        description: Список тегов с пользовательским разделителем
        options:
          go_tags: envSeparator:";"  # Использовать ; в качестве разделителя
        example: "tag1;tag2;tag3"

      - name: privateLabels
        type: "map[string]string"
        description: Пары ключ-значение с пользовательскими разделителями
        options:
          go_tags: envSeparator:";" envKeyValSeparator:"="  # Разделители для списка и пар ключ-значение
        example: "key1=value1;key2=value2"
  - name: Config
    fields:
      - name: ConfigPath
        type: string
        description: Путь к конфигурационному файлу
        required: true
        options:
          go_env_options: file  # Проверить существование файла
        example: "/etc/app/config.json"
```

Результат выполнения:

```go
// Webserver
type Webserver struct {
  Config
  ApiKey        string            `env:"API_KEY,required,unset,notEmpty"`                        // API ключ, который будет очищен после чтения
  Tags          []string          `env:"TAGS" envSeparator:";"`                                  // Список тегов с пользовательским разделителем
  privateLabels map[string]string `env:"PRIVATE_LABELS" envSeparator:";" envKeyValSeparator:"="` // Пары ключ-значение с пользовательскими разделителями
}

// Config
type Config struct {
    ConfigPath string `env:"CONFIG_PATH,required,file"` // Путь к конфигурационному файлу
}
```

### Опции Markdown шаблона

Шаблон `markdown` поддерживает глобальные опции:

```yaml
options:
  # Основные настройки markdown
  md_title: Заголовок Markdown файла                # Заголовок документации
  md_description: |                                 # Дополнительное описание в начале файла
    Дополнительное описание вверху страницы
  
  # Настройки для раздела типов
  md_types_title: Заголовок для типов              # Заголовок раздела с типами
  md_types_description: |                          # Дополнительное описание для раздела типов
    Дополнительное описание для типов

  # Скрытие столбцов в таблице групп
  md_groups_hide_type: true        # Скрыть столбец с типом переменной
  md_groups_hide_required: true    # Скрыть столбец с обязательностью поля
  md_groups_hide_default: true     # Скрыть столбец со значением по умолчанию
  md_groups_hide_example: true     # Скрыть столбец с примером значения
  md_groups_hide_description: true # Скрыть столбец с описанием

  # Скрытие столбцов в таблице типов
  md_types_hide_type: true        # Скрыть столбец с типом
  md_types_hide_import: true      # Скрыть столбец с путем импорта
  md_types_hide_description: true # Скрыть столбец с описанием типа
  md_types_hide_values: true      # Скрыть столбец с возможными значениями


groups:
  - ...
    options:
      md_description: "Описание для группы"
    fields:
      - ...
        options:
          md_hide: true  # Скрыть поле
```

По умолчанию все столбцы отображаются. Чтобы скрыть определенный столбец, установите соответствующую опцию в `true`. Столбец `Name` в таблице типов всегда отображается.

### Опции Example шаблона

Шаблон не использует каких-либо специальных опций.

## Разработка

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
  - name: CustomField
    type: CustomType
```

### Какие функции доступны в шаблонах?

В шаблонах доступны следующие встроенные функции:

- Функции для работы со строками:
  - `title` - преобразование первой буквы в верхний регистр
  - `upper` - преобразование в верхний регистр
  - `lower` - преобразование в нижний регистр
  - `camel` - преобразование в camelCase
  - `snake` - преобразование в snake_case
  - `kebab` - преобразование в kebab-case
  - `pascal` - преобразование в PascalCase
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
  - `oneline` - преобразование многострочного текста в одну строку
  - `isURL` - проверка является ли строка URL

- Функции для работы с типами:
  - `toString` - преобразование в строку
  - `toInt` - преобразование в целое число
  - `toBool` - преобразование в логическое значение
  - `findType` - поиск информации о типе
  - `getImports` - получение списка импортов

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
  - `processTemplate` - обработка шаблона с использованием доступных функций

- Функции для работы с путями:
  - `pathDir` - получение имени директории из пути
  - `pathBase` - получение имени файла из пути
  - `pathExt` - получение расширения файла
  - `pathRel` - получение относительного пути
  - `getConfigPath` - путь к файлу конфигурации
  - `getOutputPath` - путь к выходному файлу
  - `getTemplatePath` - путь к файлу шаблона

- Функции для работы с Go:
  - `goCommentGenerate` - генерация комментария go:generate

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
{{ pathBase "/path/to/file.txt" }}  // Результат: file.txt

// Работа с URL
{{ if isURL "https://github.com" }}
  // URL валиден
{{ end }}

// Обработка многострочного текста
{{ $text := "строка1\nстрока2" | oneline }}  // Результат: строка1 строка2
```

## Лицензия

MIT