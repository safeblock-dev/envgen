# Application Environment Configuration

This document describes all environment variables used in the application.
Please ensure all required variables are properly set before running the application.
For local development, you can use the provided `.env.example` file as a template.


## App

Main application settings

| Name | Type | Required | Default | Example | Description |
|--------|------|----------|---------|---------|-------------|
| `MODE` | [`Environment`](#custom-types) (string) | ✗ | `production` | `production` | Application running mode (Possible values: development, staging, production) |
| `DEBUG` | bool | ✗ | `false` | - | Enable debug logging |
| `PORT` | int | ✓ | - | `8080` | HTTP server port |

## Database

Database connection settings

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