# Application Environment Configuration

This document describes all environment variables used in the application.
Please ensure all required variables are properly set before running the application.
For local development, you can use the provided `.env.example` file as a template.


## App

Main application settings

Core application configuration settings.

> **Important**: Changes to these settings require application restart.

For development mode, set `APP_MODE=development` and `APP_DEBUG=true`.

See [Database](#database) section for database connection configuration.


| Name | Type | Required | Default | Example | Description |
|--------|------|----------|---------|---------|-------------|
| `MODE` | [`Environment`](#custom-types) | ✗ | `production` | `production` | Application running mode (Possible values: development, staging, production) |
| `DEBUG` | bool | ✗ | `false` | - | Enable debug logging |
| `PORT` | int | ✓ | - | `8080` | HTTP server port |

## Database

Database connection settings

PostgreSQL database connection configuration.

Connection string format:
```
postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}
```

For local development, default connection string would be:
```
postgresql://postgres:postgres@localhost:5432/app
```


| Name | Type | Required | Default | Example | Description |
|--------|------|----------|---------|---------|-------------|
| `DB_HOST` | string | ✗ | `localhost` | - | Database host address |
| `DB_PORT` | int | ✗ | `5432` | - | Database port number |
| `DB_NAME` | string | ✓ | - | - | Database name |
| `DB_USER` | string | ✓ | - | - | Database user |
| `DB_PASSWORD` | string | ✓ | - | - | Database password |

## Custom Environment Types

The following section describes custom types used in the configuration.
These types provide additional validation and documentation for specific environment variables.
Each type includes its possible values.


| Name | Type | Import Path | Description | Possible Values |
|----|------|------------|-------------|----------------|
| `Environment` | string | - | Application environment | `development`, `staging`, `production` | 