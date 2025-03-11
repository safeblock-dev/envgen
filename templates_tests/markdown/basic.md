# Environment Variables Documentation

## App

Basic application settings

| Variable | Type | Required | Default | Example | Description |
|----------|------|----------|---------|---------|-------------|
| `DEBUG` | bool | ✗ | `false` | - | Enable debug mode |
| `PORT` | int | ✓ | - | `8080` | Server port |
| `HOST` | string | ✗ | `localhost` | - | Server host |

## Database

Database connection settings

| Variable | Type | Required | Default | Example | Description |
|----------|------|----------|---------|---------|-------------|
| `DB_URL` | string | ✓ | - | `postgres://user:pass@localhost:5432/db` | Database connection URL |
| `DB_LOG_LEVEL` | string | ✗ | `info` | - | Database logging level (Possible values: debug, info, warn, error) |

## Types

### LogLevel

Log level for the application

Possible values:
- `debug`
- `info`
- `warn`
- `error` 