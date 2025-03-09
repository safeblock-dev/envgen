# envgen

`envgen` is a flexible tool for generating environment configuration code. It creates type-safe configuration structures and documentation from YAML definitions.

## Features

- Multiple output formats:
  - Go structures with `env` tags
  - Environment file templates
  - Markdown documentation
- Type system with custom type definitions
- Group-based configuration organization
- Environment variable prefix support
- Customizable templates
- Type and group filtering for composite configurations
- Extended validation and error reporting

## Installation

```bash
go install github.com/safeblock-dev/envgen/cmd/envgen@latest
```

## Quick Start

1. Create a configuration file `config.yaml`:
```yaml
# Example configuration for Go package generation
options:
  go_package: myapp/config  # Package name for generated code

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
        required: true
        example: "0.0.0.0"
      
      - name: env
        type: Environment
        description: Environment
        default: "development"
        required: true
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
# Generate using local template
envgen -c config.yaml -o config.go -t ./templates/config.tmpl

# Generate using template from URL
envgen --config config.yaml --out config.go --template https://raw.githubusercontent.com/user/repo/template.tmpl

# Generate ignoring specific types and groups
envgen -c config.yaml -o config.go -t ./templates/config.tmpl --ignore-types Duration,URL --ignore-groups Database

# Show version
envgen version
```

This will create a `config.go` file with type-safe structures for configuration management. The generated code will use environment variables with the `SERVER_` prefix (e.g., `SERVER_PORT`, `SERVER_HOST`, `SERVER_ENV`).

## Configuration Format

### Types

Types allow defining custom types with validation and documentation:

```yaml
types:
  - name: Duration      # Required: Type name for referencing in fields
    type: time.Duration # Required: Type definition (built-in or custom)
    import: time       # Optional: Import path for custom types
    description: Time interval # Optional: Type description
    values:           # Optional: Possible values for documentation
      - 1s
      - 1m
```

### Groups

Groups organize related configuration fields:

```yaml
groups:
  - name: Database     # Required: Group name
    description: Database settings # Optional: Group description
    prefix: DB_         # Optional: Environment variable prefix
    options:            # Optional: Group-specific options
      go_name: DBConfig # Optional: Override struct name
    fields:             # Required: At least one field must be defined
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
  - name: port        # Required: Environment variable name
    type: int        # Required: Field type (built-in or custom)
    description: Port # Optional: Field description
    default: "8080"  # Optional: Default value
    required: true   # Optional: Whether the field is required
    example: "9000"  # Optional: Example value for documentation
    options:         # Optional: Field-specific options
      import: "custom/pkg" # Optional: Import path for custom types
      name_field: Port    # Optional: Override struct field name
```

## Advanced Features

### Composite Configurations

You can use groups as types and filter them during generation:

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

### Using Generated Code

After generating the code, you can use it in your Go application:

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

    fmt.Printf("Server will start on %s:%d\n", cfg.Host, cfg.Port)
    fmt.Printf("Environment: %s\n", cfg.Env)
}
```

The generated code uses the `env` package for parsing environment variables. Make sure to add it to your dependencies:

```bash
go get github.com/caarlos0/env/v11
```

### Go Template Options

The `go-env` template supports additional options for fields that control how environment variables are processed:

```yaml
fields:
  - name: api_key
    type: string
    description: API key from file
    required: true      # Global option: adds ,required tag
    options:
      go_file: true     # Read value from file
      go_expand: true   # Enable environment variable expansion
      go_init: true     # Initialize nil pointers
      go_notEmpty: true # Error if value is empty
      go_unset: true    # Unset environment variable after reading
```

#### Available Go Options

- `go_file`: Indicates that the value should be read from the file specified by the environment variable
- `go_expand`: Enables environment variable expansion in values (e.g., `FOO_${BAR}`)
- `go_init`: Initializes pointers that would otherwise be nil
- `go_notEmpty`: Returns an error if the environment variable is empty
- `go_unset`: Unsets environment variables after they are read

These options are implemented using the `github.com/caarlos0/env/v11` package tags.

## Development

### Project Structure

```
.
├── cmd/envgen/          # CLI application
├── pkg/envgen/          # Core package
│   ├── config.go        # Configuration types and validation
│   ├── envgen.go        # Main generation logic
│   ├── template.go      # Template loading
│   ├── template_context.go # Template context and functions
│   └── templatefuncs/   # Template helper functions
├── templates/           # Built-in templates
└── templates_tests/     # Tests and examples
```

### Running Tests

```bash
go test ./...
```

Update golden files for template tests:
```bash
UPDATE_GOLDEN=1 go test ./templates_tests
```

## License

MIT