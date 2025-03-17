# Environment Variables Documentation

## Server

Server settings

| Name | Type | Required | Default | Example | Description |
|--------|------|----------|---------|---------|-------------|
| `SERVER_PORT` | int | ✓ | `8080` | - | Server port |
| `SERVER_ENV` | [`Environment`](#custom-types) | ✓ | `development` | - | Environment (Possible values: development, staging, production) |

## Custom Types

| Name | Type | Import Path | Description | Possible Values |
|----|------|------------|-------------|----------------|
| `Environment` | string | - | Application environment | `development`, `staging`, `production` | 