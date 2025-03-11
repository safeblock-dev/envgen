# Environment Variables Documentation

## App

Basic application settings

| Name | Type | Required | Default | Example | Description |
|--------|------|----------|---------|---------|-------------|
| `DEBUG` | bool | ✗ | `false` | - | Enable debug mode |
| `PORT` | int | ✓ | - | `8080` | Server port |
| `HOST` | string | ✗ | `localhost` | - | Server host |

## Database

Database connection settings

| Name | Type | Required | Default | Example | Description |
|--------|------|----------|---------|---------|-------------|
| `DB_URL` | string | ✓ | - | `postgres://user:pass@localhost:5432/db` | Database connection URL |
| `DB_LOG_LEVEL` | string | ✗ | `info` | - | Database logging level (Possible values: debug, info, warn, error) |

## Custom Types

| Type | Import Path | Description | Possible Values |
|----|------------|-------------|----------------|
| `LogLevel` | - | Log level for the application | `debug`, `info`, `warn`, `error` | 