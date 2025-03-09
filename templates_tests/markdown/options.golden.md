# Application Configuration Guide

This document describes all environment variables that can be used to configure the application.
Each variable is documented with its type, whether it's required, and its default value if any.


## App

Main application settings

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `MODE` | string | ✗ | `production` | Application running mode |
| `DEBUG` | bool | ✗ | `false` | Enable debug logging |
| `PORT` | int | ✓ | - | HTTP server port |

## Database

Database connection settings

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `DB_HOST` | string | ✗ | `localhost` | Database host address |
| `DB_PORT` | int | ✗ | `5432` | Database port number |
| `DB_NAME` | string | ✓ | - | Database name |
| `DB_USER` | string | ✓ | - | Database user |
| `DB_PASSWORD` | string | ✓ | - | Database password | 