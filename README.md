# envgen

**`envgen` – A powerful tool for generating environment configuration from YAML**

Creates type-safe Go structures, `.env.example` files, documentation, and any other files using custom templates.

[Russian Documentation (Русская версия)](README.ru.md)

### Advantages

- 🔒 **Type Safety**: Type validation at compile time
- 🔄 **Automatic Generation**: Documentation and examples are always synchronized with code
- 🎨 **Any Template**: Support for custom templates and formats
- 🛠 **Flexible Configuration**: Simple settings for customization
- 📝 **Auto-documentation**: Automatically generated Markdown documentation
- 🔍 **Transparency**: Clear configuration structure in YAML format

## Features

- Multiple output formats:
  - Go structures with `env` tags
  - Environment file examples (`.env.example`)
  - Documentation in Markdown format
- Configurable templates
- Support for custom templates

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
    prefix: SERVER_  # Prefix for environment variables
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
- `-t, --template`: Path to template file or URL (required)
- `--ignore-types`: Comma-separated list of types to ignore
- `--ignore-groups`: Comma-separated list of groups to ignore

Examples:
```bash
# Generation using local template
envgen -c config.yaml -o config.go -t ./templates/config.tmpl

# Generation using template from URL
envgen --config config.yaml --out config.go --template https://raw.githubusercontent.com/user/repo/template.tmpl

# Generation ignoring specific types and groups
envgen -c config.yaml -o config.go -t ./templates/config.tmpl --ignore-types Duration,URL --ignore-groups Database

# Show version
envgen version
```

This will create a `config.go` file with type-safe structures for configuration handling. The generated code will use environment variables with the `SERVER_` prefix (e.g., `SERVER_PORT`, `SERVER_HOST`, `SERVER_ENV`).

## Configuration Format

### Options

Options allow you to customize and modify information in the template. Different templates use different options.

```yaml
options:
  go_package: mypkg
```

This example shows an option that sets the package name for the `go-env` template. You can also set options in a group and in individual fields.

### Groups

Groups organize related configuration fields:

```yaml
groups:
  - name: Database     # Required: group name
    description: Database settings # Optional: group description
    prefix: DB_         # Optional: prefix for environment variables
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

To use created types, specify their name as the type value in the field description:

```yaml
fields:
  - name: github                  
    type: AppURL                # Specify type name
    example: "http://github.com/safeblock-dev" 
  - name: twitter                  
    type: AppURL                # Type can be used multiple times
    example: "http://x.com/safeblock" 
```

The following special keys are available in options:
- `{{ ConfigPath }}` - outputs the path to the configuration file
- `{{ OutputPath }}` - outputs the path to the output file
- `{{ TemplatePath }}` - outputs the path to the template

Special options are also available for groups and fields:

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

- `go-env`: Generates Go structures with `env` tags
- `example`: Creates `.env` file templates
- `markdown`: Creates documentation in Markdown format

Examples of using built-in templates:

```bash
# Generate Go structures
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
    # Generate configuration
    //go:generate envgen -c {{ ConfigPath }} -o {{ OutputPath }} -t {{ TemplatePath }}
    # Generate documentation
    //go:generate envgen -c {{ ConfigPath }} -o docs/{{ OutputPath }} -t markdown
  go_meta: |
    // Version: v0.1.2
    // Template: {{ TemplatePath }}
```

The `go_package` option is required for the `go-env` template. If not specified, `envgen` will try to use the folder name from the `out` flag, but this is considered bad practice since if the path is something like `config.go`, the package name will be set to `.`.

The `go_generate` option lets you specify custom code generation commands. If not specified, the default command is used.

The `go_meta` option adds additional information after the `go_generate` block.

The following special keys are available in options:
- `{{ ConfigPath }}` - outputs the path to the configuration file
- `{{ OutputPath }}` - outputs the path to the output file
- `{{ TemplatePath }}` - outputs the path to the template

Special options are also available for groups and fields:

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

Result:

```go
// App name changed to CustomAppConfig
type CustomAppConfig struct {
   // DebugMode name (debug_mode in config file) changed to IsDebug
	IsDebug bool `env:"DEBUG_MODE" envDefault:"false"`
}
```

### Markdown Template Options

The `markdown` template supports global options:

```yaml
options:
  md_title: Markdown File Title
  md_description: |
    Any additional description
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

### How to Add a Custom Template?

Create a template file with `.tmpl` or `.tpl` extension. Use Go templates syntax and available context functions. Example of a simple template:

```go
// File: custom.tmpl
{{- range $group := .Groups }}
# {{ $group.Description }}
{{- range $field := $group.Fields }}
{{ $field.Name | upper }}={{ $field.Default }}  # {{ $field.Description }}
{{- end }}

{{- end }}
```

Generate using template:
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

### How to Use Custom Types?

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

### What Functions are Available in Templates?

The following built-in functions are available in templates:

- String manipulation functions:
  - `title` - convert first letter to uppercase
  - `upper` - convert to uppercase
  - `lower` - convert to lowercase
  - `pascal` - convert to PascalCase
  - `camel` - convert to camelCase
  - `snake` - convert to snake_case
  - `kebab` - convert to kebab-case
  - `append` - append string to end
  - `uniq` - remove duplicates
  - `slice` - get substring
  - `contains` - check for substring
  - `hasPrefix` - check prefix
  - `hasSuffix` - check suffix
  - `replace` - replace substring
  - `trim` - remove whitespace
  - `join` - join strings
  - `split` - split string

- Type manipulation functions:
  - `toString` - convert to string
  - `toInt` - convert to integer
  - `toBool` - convert to boolean
  - `findType` - find type information
  - `getImports` - get import list
  - `typeImport` - get type import

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
  - `getConfigPath` - path to configuration file
  - `getOutputPath` - path to output file
  - `getTemplatePath` - path to template file

Usage example:
```go
{{ $name := "my_variable" }}
{{ $name | pascal }}  // Result: MyVariable
{{ $name | upper }}   // Result: MY_VARIABLE

{{ if hasOption "go_package" }}
package {{ getOption "go_package" }}
{{ end }}

// Date and time operations
{{ datetime }}  // Result: 2024-03-21 15:04:05

// Conditional operations
{{ $value := "test" | default "default_value" }}

// Path operations
{{ $dir := getFileName "/path/to/file.txt" }}  // Result: file.txt
```

## License

MIT 