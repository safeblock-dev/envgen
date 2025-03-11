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
      - name: port
        type: int
        description: Server port
        default: "8080"
        required: true
        example: "9000"
      
      - name: host
        type: string
        description: Server host
        default: "localhost"
        example: "0.0.0.0"
      
      - name: env
        type: Environment
        description: Environment
        default: "development"
        example: "production"
```

2. Generate Go code:
```bash
envgen -c config.yaml -o config.go -t go-env
```

### Command Line Flags

The tool supports the following flags:

- `-c, --config`: Path to input YAML configuration file (required)
- `-o, --out`: Path to output file (required)
- `-t, --template`: Path to template or URL (required)
- `--ignore-types`: Comma-separated list of types to ignore
- `--ignore-groups`: Comma-separated list of groups to ignore

Examples:
```bash
# Generate using local template
envgen -c config.yaml -o config.go -t ./templates/config.tmpl

# Generate using template from URL
envgen --config config.yaml --out config.go --template https://raw.githubusercontent.com/user/repo/template.tmpl

# Generate database configurations only
envgen -c config.yaml -o config.go -t ./templates/config.tmpl --ignore-types Duration,URL --ignore-groups Database

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
      - name: host
        type: string
        description: Database host
        required: true
        default: localhost
```

### Fields

Fields represent individual environment variables:

```yaml
fields:
  - name: url                   # Required: environment variable name
    type: string                # Required: field type (built-in or custom)
    description: API endpoint   # Optional: field description
    default: "http://127.0.0.1" # Optional: default value
    required: true              # Optional: whether the field is required
    example: "http://test.com"  # Optional: example value for documentation
    options:                    # Optional: additional field parameters
      go_name: "URL"            # Optional: any template option
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
  - name: github                  
    type: AppURL                # Specify type name
    example: "http://github.com/safeblock-dev" 
  - name: twitter                  
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
      - name: host
        type: string
        default: localhost

  - name: Redis
    description: Redis settings
    prefix: REDIS_
    fields:
      - name: port
        type: int
        default: "6379"

  - name: Webserver
    description: Web server configuration
    fields:
      - name: db
        type: Postgres
        options:
          name_field: DB
      - name: cache
        type: Redis
```

Generate only database configurations:
```bash
envgen -c config.yaml -o config.go -t go-env --ignore-groups Webserver
```

This is especially useful when you have structures that you don't want to show, for example, in `.env.example`.

### Templates

The tool includes three built-in templates:

- `go-env`: Generates Go structs with `env` tags
- `example`: Creates `.env` file templates
- `markdown`: Creates Markdown documentation

Examples of using built-in templates:

```bash
# Generate Go structs
envgen -c config.yaml -o config.go -t go-env

# Generate .env file template
envgen -c config.yaml -o .env.example -t example

# Generate documentation
envgen -c config.yaml -o config.md -t markdown
```

### Go Template Options

The `go-env` template supports global options:

```yaml
options:
  go_package: config # Required field
  go_generate: |
    # Configuration generation
    //go:generate envgen -c {{ ConfigPath }} -o {{ OutputPath }} -t {{ TemplatePath }}
    # Documentation generation
    //go:generate envgen -c {{ ConfigPath }} -o docs/{{ OutputPath }} -t markdown
  go_meta: |
    // Version: v0.1.2
    // Template: {{ TemplatePath }}
```

The `go_package` option is required for the `go-env` template. If no value is specified, `envgen` will try to use the folder name from the `out` flag, but this is considered bad practice since if the path is like `config.go`, the package name will be set as `.`, which will lead to a compilation error.

The `go_generate` option allows you to specify custom code generation commands. If this option is not specified, the default command is used.

The `go_meta` option will add additional information after the `go_generate` block.

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
      - name: debug_mode
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
      - name: sentry
        type: SentryConfig
      - name: grpc_port
        type: int
        default: "8002"
      - name: http_port
        type: int
        default: "8001"
        options:
          go_name: HTTP_PORT
          go_tags: env:"NOT_SKIPPED"
    
  - name: CustomEnvTags
    description: Selective application of env tags
    fields:
      - name: not_skipped
        type: string
      - name: debug
        type: bool
        options:
          go_skip_env_tag: true
      - name: port
        type: int
        options:
          go_skip_env_tag: true
          go_env_options: skipped
          go_tags: env:"NOT_SKIPPED,required,notEmpty"
```

Field-specific options:

- `go_env_options` - allows adding additional options to the `env` tag. For example: `file`, `unset`, `notEmpty`, and other options. All options are passed directly to the tags without additional validation.
- `go_tags` - allows adding additional tags to the structure. Supports specifying any tags without restrictions. When used with the [`env`](github.com/caarlos0/env/v11) package, commonly used options include:
  - `envSeparator` - separator for slices
  - `envKeyValSeparator` - separator for key-value pairs in maps

Example usage:

```yaml
fields:
  - name: config_path
    type: string
    description: Path to configuration file
    required: true
    options:
      go_env_options: file  # Check if file exists
    example: "/etc/app/config.json"
  
  - name: api_key
    type: string
    description: API key that will be cleared after reading
    required: true
    options:
      go_env_options: unset,notEmpty  # Clear after reading and check for emptiness
    example: "secret-key"
  
  - name: tags
    type: "[]string"
    description: List of tags with custom separator
    options:
      go_tags: envSeparator:";"  # Use ; as separator
    example: "tag1;tag2;tag3"

  - name: labels
    type: "map[string]string"
    description: Key-value pairs with custom separators
    options:
      go_tags: envSeparator:";" envKeyValSeparator:"="  # Separators for list and key-value pairs
    example: "key1=value1;key2=value2"
```

Execution result:

```go
type AppConfig struct {
	ConfigPath string `env:"APP_CONFIG_PATH,required,file"` // Path to configuration file
	ApiKey string `env:"APP_API_KEY,required,unset,notEmpty"` // API key that will be cleared after reading
	Tags []string `env:"APP_TAGS" envSeparator:";"` // List of tags with custom separator
	Labels map[string]string `env:"APP_LABELS" envSeparator:";" envKeyValSeparator:"="` // Key-values with custom separators
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
    Additional descript

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
```

These options are additional and not required.

### Example Template Options

The template does not use any special options.

## Development

### Project Structure

```
.
├── cmd/envgen/             # CLI application
├── pkg/envgen/             # Main package
│   ├── config.go           # Configuration types and validation
│   ├── envgen.go           # Main generation logic
│   ├── template.go         # Template loading
│   ├── template_context.go # Template context and functions
│   └── templatefuncs/      # Template helper functions
├── templates/              # Built-in templates
└── templates_tests/        # Tests and examples
```

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
  - name: custom_field
    type: CustomType
```

### What functions are available in templates?

The following built-in functions are available in templates:

- String manipulation functions:
  - `title` - converts first letter to uppercase
  - `upper` - converts to uppercase
  - `lower` - converts to lowercase
  - `pascal` - converts to PascalCase
  - `camel` - converts to camelCase
  - `snake` - converts to snake_case
  - `kebab` - converts to kebab-case
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

- Type manipulation functions:
  - `toString` - converts to string
  - `toInt` - converts to integer
  - `toBool` - converts to boolean
  - `findType` - finds type information
  - `getImports` - gets import list
  - `typeImport` - gets type import

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

- Path manipulation functions:
  - `getDirName` - get directory name
  - `getFileName` - get file name
  - `getFileExt` - get file extension
  - `joinPaths` - join paths
  - `getConfigPath` - configuration file path
  - `getOutputPath` - output file path
  - `getTemplatePath` - template file path

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
{{ $dir := getFileName "/path/to/file.txt" }}  // Result: file.txt
```

## License

MIT