# envgen

**`envgen` – A powerful tool for generating environment configuration from YAML**

Creates type-safe Go structs, `.env.example`, documentation, and any other files using custom templates

[Russian Documentation (Русская версия)](README.ru.md)

### Benefits

- 🔒 **Type Safety**: Type validation at compile time
- 🔄 **Automatic Generation**: Documentation and examples always in sync with code
- 🎨 **Any Template**: Support for custom templates and formats
- 🛠 **Flexible Configuration**: Use simple settings for customization
- 📝 **Auto-documentation**: Markdown documentation generated automatically
- 🔍 **Transparency**: Clear configuration structure in YAML format

## Features

- Various output formats:
  - Go structs with `env` tags
  - Environment example files (`.env.example`)
  - Markdown documentation
- Customizable templates
- Ability to use custom templates

## Installation

```bash
go install github.com/safeblock-dev/envgen/cmd/envgen@latest
```

## Quick Start

1. Create a configuration file `config.yaml`:
```yaml
# Example configuration for Go package generation
options:
  go_package: config  # Package name for generated code

types:
  - name: Environment
    type: string
    description: Application environment (development, staging, production)
    values:
      - development
      - staging
      - production

groups:
  - name: Server
    description: Web server settings
    prefix: SERVER_  # Environment variable prefix
    fields:
      - name: Port
        type: int
        description: Server port
        default: "8080"
        required: true
        example: "9000"
      
      - name: Host
        type: string
        description: Server host
        default: "localhost"
        example: "0.0.0.0"
      
      - name: ENV
        type: Environment
        description: Environment
        default: "development"
        example: "production"
```

2. Generate Go code:
```bash
envgen gen -c config.yaml -o config.go -t go-env
```

### Commands

`envgen` supports the following commands:

- `gen` (or `generate`): Generate configuration files
  - `-c, --config`: Path to input YAML configuration file (required)
  - `-o, --out`: Path to output file (required)
  - `-t, --template`: Path to template or URL (required)
  - `--ignore-types`: Comma-separated list of types to ignore
  - `--ignore-groups`: Comma-separated list of groups to ignore

- `ls` (or `templates`, `list`): List available standard templates

- `version`: Show program version

Examples:
```bash
# Generate using local template
envgen gen -c config.yaml -o config.go -t ./templates/config.tmpl

# Generate using template from URL
envgen gen --config config.yaml --out config.go --template https://raw.githubusercontent.com/user/repo/template.tmpl

# Generate with ignoring specific types and groups
envgen gen -c config.yaml -o config.go -t ./templates/config.tmpl --ignore-types Duration,URL --ignore-groups Database

# Show available templates
envgen ls

# Show version
envgen version
```

This will create a `config.go` file with type-safe structs for configuration handling. The generated code will use environment variables with the `SERVER_` prefix (e.g., `SERVER_PORT`, `SERVER_HOST`, `SERVER_ENV`).

## Configuration Format

### Options

Options enable you to configure and modify information in the template. Different templates use different options.

```yaml
options:
  go_package: mypkg
```

This example shows an option that will set the package name for the `go-env` template. You can also set options in a group and in an individual field.

### Groups

Groups organize related configuration fields:

```yaml
groups:
  - name: Database     # Required: group name
    description: Database settings # Optional: group description
    prefix: DB_         # Optional: environment variable prefix
    options:            # Optional: group parameters
      go_name: DBConfig # Optional: any template option
    fields:             # Required: at least one field must be defined
      - name: Host
        type: string
        description: Database host
        required: true
        default: localhost
```

### Fields

Fields represent individual environment variables:

```yaml
fields:
  - name: URL                   # Required: environment variable name
    type: string                # Required: field type (built-in or custom)
    description: API endpoint   # Optional: field description
    default: "http://127.0.0.1" # Optional: default value
    required: true              # Optional: whether the field is required
    example: "http://test.com"  # Optional: example value for documentation
    options:                    # Optional: additional field parameters
      go_name: "GitURL"         # Optional: any template option
```

### Types

Types allow you to define custom types, add context to a type, and reuse them:

```yaml
types:
  - name: Duration        # Required: type name for field references
    type: time.Duration   # Required: type definition (built-in or custom)
    import: time          # Optional: import path for custom types
    description: Interval # Optional: type description
    values:               # Optional: possible values for documentation
      - 1s
      - 1m
```

You can create multiple similar types:

```yaml
types:
  - name: AppENV
    type: string
    description: Environment name
    values: ["prod", "dev"]
  - name: MediaURL
    type: string
    description: Media source link
```

To use created types, specify their `name` as the `type` value in the field description:

```yaml
fields:
  - name: Github                  
    type: AppURL                # Specify type name
    example: "http://github.com/safeblock-dev" 
  - name: Twitter                  
    type: AppURL                # Type can be used multiple times
    example: "http://x.com/safeblock" 
```

## Advanced Features

### Composite Configurations

You can ignore types and groups during generation:

```yaml
groups:
  - name: Postgres
    description: PostgreSQL settings
    prefix: PG_
    fields:
      - name: Host
        type: string
        default: localhost

  - name: Redis
    description: Redis settings
    prefix: REDIS_
    fields:
      - name: Port
        type: int
        default: "6379"

  - name: Webserver
    description: Web server configuration
    fields:
      - name: DB
        type: Postgres
      - name: Cache
        type: Redis
```

Generate only database configurations:
```bash
envgen gen -c config.yaml -o config.go -t go-env --ignore-groups Webserver
```

This is especially useful when you have structures that you don't want to show, for example, in `.env.example`.

### Templates

The tool includes four built-in templates:

- `go-env`: Generates Go structs with `env` tags
- `go-env-example`: Creates `.env.example` templates considering `go-env` tags (options)
- `example`: Creates `.env.example` templates
- `markdown`: Creates Markdown documentation

Examples of using built-in templates:

```bash
# Generate Go structs
envgen gen -c config.yaml -o config.go -t go-env

# Generate .env file template
envgen gen -c config.yaml -o .env.example -t example

# Generate documentation
envgen gen -c config.yaml -o config.md -t markdown
```

### Go Template Options

The `go-env` template supports global options:

```yaml
options:
  go_package: config # Optional field
  go_meta: |
    # Configuration generation
    {{ goCommentGenerate "" "" "" }}
    # Documentation generation
    {{ goCommentGenerate "" "docs/README.md" "../../templates/markdown" }}
    // Version: v0.1.2
```

If the `go_package` value is not specified, `envgen` will attempt to use the folder name from the `out` flag.

The `go_meta` option allows you to specify custom commands for code generation. If this option is not specified, the default command is used. If you don't want the `//go:generate` output, leave the `go_meta` field empty.

The `go_meta` option allows you to call any template functions from the [funcs.go](pkg/envgen/funcs.go) file (for example `title`, `upper`, etc.).

The following special keys are available in options (`go_generate`, `go_meta`):
- `{{ ConfigPath }}` - outputs the configuration file path
- `{{ OutputPath }}` - outputs the output file path
- `{{ TemplatePath }}` - outputs the template path

Special options for customizing names for groups and fields are also available:

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

Execution result:

```go
// App name changed to CustomAppConfig
type CustomAppConfig struct {
   // DebugMode name (debug_mode in config file) changed to IsDebug
	IsDebug bool `env:"DEBUG_MODE" envDefault:"false"`
}
```

Additional options for configuring `env` tags for groups and fields:

- `go_skip_env_tag` - disables the generation of the `env` tag

Example usage:

```yaml
  - name: NoEnvTag
    description: Group without env tags
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
    description: Selective application of env tags
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

Field-specific options:

- `go_include` - if true, uses Go struct embedding
- `go_env_options` - allows adding additional options to the `env` tag. For example: `file`, `unset`, `notEmpty`, and other options. All options are passed directly to the tags without additional validation.
- `go_tags` - allows adding additional tags to the structure. Supports specifying any tags without restrictions. When used with the [`env`](github.com/caarlos0/env/v11) package, commonly used options include:
  - `envSeparator` - separator for slices
  - `envKeyValSeparator` - separator for key-value pairs in maps

Example usage:

```yaml
groups:
  - name: Webserver
    fields:
      - name: Config
        type: Config
        options:
          go_skip_env_tag: true
          go_include: true # Struct embedding

      - name: ApiKey
        type: string
        description: API key that will be cleared after reading
        required: true
        options:
          go_env_options: unset,notEmpty  # Clear after reading and check for emptiness
        example: "secret-key"

      - name: Tags
        type: "[]string"
        description: List of tags with custom separator
        options:
          go_tags: envSeparator:";"  # Use ; as separator
        example: "tag1;tag2;tag3"

      - name: privateLabels
        type: "map[string]string"
        description: Key-value pairs with custom separators
        options:
          go_tags: envSeparator:";" envKeyValSeparator:"="  # Separators for list and key-value pairs
        example: "key1=value1;key2=value2"
  - name: Config
    fields:
      - name: ConfigPath
        type: string
        description: Path to configuration file
        required: true
        options:
          go_env_options: file  # Check if file exists
        example: "/etc/app/config.json"
```

Execution result:

```go
// Webserver
type Webserver struct {
  Config
  ApiKey        string            `env:"API_KEY,required,unset,notEmpty"`                        // API key that will be cleared after reading
  Tags          []string          `env:"TAGS" envSeparator:";"`                                  // List of tags with custom separator
  privateLabels map[string]string `env:"PRIVATE_LABELS" envSeparator:";" envKeyValSeparator:"="` // Key-value pairs with custom separators
}

// Config
type Config struct {
    ConfigPath string `env:"CONFIG_PATH,required,file"` // Path to configuration file
}
```

### Markdown Template Options

The `markdown` template supports global options:

```yaml
options:
  md_title: Markdown File Title
  md_description: |
    Additional description at the top of the page
  md_types_title: Types Section Title
  md_types_description: |
    Additional description

  # Hide specific columns in the groups table
  md_groups_hide_type: true        # Hide Type column
  md_groups_hide_required: true    # Hide Required column
  md_groups_hide_default: true     # Hide Default column
  md_groups_hide_example: true     # Hide Example column
  md_groups_hide_description: true # Hide Description column

  # Hide specific columns in the types table
  md_types_hide_type: true         # Hide Type column
  md_types_hide_import: true       # Hide Import column
  md_types_hide_description: true  # Hide Description column
  md_types_hide_values: true       # Hide Possible Values column

groups:
  - ...
    options:
      md_description: Additional description
    fields:
      - ...
        options:
          md_hide: true  # This field will be hidden in documentation 
```

These options are additional and not required.

### Example Template Options

The template does not use any special options.

## Development

### Running Tests

```bash
go test ./...
```

Update golden files for template tests:
```bash
UPDATE_GOLDEN=1 go test ./templates_tests
```

## Frequently Asked Questions

### How to add your own template?

Create a template file with `.tmpl` or `.tpl` extension. Use Go templates syntax and available functions from context. Simple template example:

```go
// File: custom.tmpl
{{- range $group := .Groups }}
# {{ $group.Description }}
{{- range $field := $group.Fields }}
{{ $field.Name | upper }}={{ $field.Default }}  # {{ $field.Description }}
{{- end }}

{{- end }}
```

Generate template:
```bash
envgen -c config.yaml -o custom.txt -t ./custom.tmpl
```

The result will look like this:
```ini
# Web server settings
PORT=8080  # Server port
HOST=localhost  # Server host
ENV=development  # Environment

# Database settings
DB_HOST=localhost  # Database host
DB_PORT=5432  # Database port
```

### How to use custom types?

1. Define type in configuration:
```yaml
types:
  - name: CustomType
    type: your/pkg.Type
    import: your/pkg
```

2. Use it in fields:
```yaml
fields:
  - name: CustomField
    type: CustomType
```

### What functions are available in templates?

The following built-in functions are available in templates:

- String manipulation functions:
  - `title` - converts first letter to uppercase
  - `upper` - converts to uppercase
  - `lower` - converts to lowercase
  - `camel` - converts to camelCase
  - `snake` - converts to snake_case
  - `kebab` - converts to kebab-case
  - `pascal` - converts to PascalCase
  - `append` - appends string to end
  - `uniq` - removes duplicates
  - `slice` - gets substring
  - `contains` - checks for substring
  - `hasPrefix` - checks prefix
  - `hasSuffix` - checks suffix
  - `replace` - replaces substring
  - `trim` - removes whitespace
  - `join` - joins strings
  - `split` - splits string
  - `oneline` - converts multiline text to single line
  - `isURL` - checks if string is a URL

- Type manipulation functions:
  - `toString` - converts to string
  - `toInt` - converts to integer
  - `toBool` - converts to boolean
  - `findType` - finds type information
  - `getImports` - gets import list

- Date and time functions:
  - `now` - current time
  - `formatTime` - format time
  - `date` - current date (YYYY-MM-DD)
  - `datetime` - current date and time (YYYY-MM-DD HH:MM:SS)

- Conditional operations:
  - `default` - default value
  - `coalesce` - first non-empty value
  - `ternary` - ternary operator
  - `hasOption` - check option existence
  - `hasGroupOption` - check group option existence
  - `getOption` - get option value
  - `getGroupOption` - get group option value
  - `processTemplate` - process template using available functions

- Path manipulation functions:
  - `pathDir` - get directory name from path
  - `pathBase` - get file name from path
  - `pathExt` - get file extension
  - `pathRel` - get relative path
  - `getConfigPath` - configuration file path
  - `getOutputPath` - output file path
  - `getTemplatePath` - template file path

- Go-specific functions:
  - `goCommentGenerate` - generate go:generate comment

Usage example:
```go
{{ $name := "my_variable" }}
{{ $name | pascal }}  // Result: MyVariable
{{ $name | upper }}   // Result: MY_VARIABLE

{{ if hasOption "go_package" }}
package {{ getOption "go_package" }}
{{ end }}

// Working with date and time
{{ datetime }}  // Result: 2024-03-21 15:04:05

// Conditional operations
{{ $value := "test" | default "default_value" }}

// Working with paths
{{ pathBase "/path/to/file.txt" }}  // Result: file.txt

// Working with URLs
{{ if isURL "https://github.com" }}
  // URL is valid
{{ end }}

// Working with multiline text
{{ $text := "line1\nline2" | oneline }}  // Result: line1 line2
```

## License

MIT