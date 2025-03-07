# envgen

`envgen` is a tool for generating environment configuration in various formats. It allows you to describe configuration in a YAML file and generate:
- Go structures with tags for [env](https://github.com/caarlos0/env)
- `.env.example` files
- Markdown documentation
- Custom formats using your own templates

## Installation

```bash
go install github.com/safeblock-dev/envgen/cmd/envgen@latest
```

## Usage

```bash
# Basic generation
envgen -c config.yaml -o config.go -t template.tmpl

# Using template from URL
envgen -c config.yaml -o config.go -t https://raw.githubusercontent.com/user/repo/template.tmpl

# Show version
envgen version

# Show help
envgen --help
```

## Built-in Templates

### go-env
Generates Go code with structures for working with environment variables:
- Support for nested structures
- Tags for [env](https://github.com/caarlos0/env)
- Required field validation
- Custom type support
- Godoc format documentation

For examples of usage, see tests in `templates_tests/go-env` directory.

### example
Creates `.env.example` file with examples of all variables:
- Default values
- Description comments
- Section grouping
- Prefix support

For examples of usage, see tests in `templates_tests/example` directory.

### markdown
Generates documentation in Markdown format:
- Complete description of all variables
- Tables with types and default values
- Usage examples
- Installation instructions
- Available in Russian (markdown-ru) and English languages

Documentation examples will be available in future releases.

Source templates:
- [templates/go-env](templates/go-env) - template for Go code
- [templates/example](templates/example) - template for .env.example
- [templates/markdown](templates/markdown) - template for English documentation
- [templates/markdown-ru](templates/markdown-ru) - template for Russian documentation

## Templates

Templates can be specified in two ways:
1. Local file path: `-t ./templates/config.tmpl`
2. URL: `-t https://raw.githubusercontent.com/user/repo/template.tmpl`

The following functions are available in templates:
- String transformations: `title`, `upper`, `lower`, `camel`, `snake`, `kebab`, `pascal`
- Type operations: `toString`, `toInt`, `toBool`
- Date and time: `now`, `formatTime`, `date`, `datetime`
- Conditional operators: `default`, `coalesce`, `ternary`
- String operations: `contains`, `hasPrefix`, `hasSuffix`, `replace`, `trim`, `join`, `split`
- Path operations: `getDirName`, `getFileName`, `getFileExt`, `joinPaths`

For detailed examples of using these functions, see tests in `pkg/envgen/templatefuncs` directory.

## Configuration

Configuration is described in a YAML file and supports:
- Variable grouping
- Group prefixes
- Description and examples for each variable
- Default values
- Required field markers
- Additional options for customization

For configuration examples, see tests in `templates_tests` directory.

## License

MIT 