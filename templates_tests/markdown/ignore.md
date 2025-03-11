# Environment Variables Documentation

## Postgres

PostgreSQL database settings

| Name | Type | Required | Default | Example | Description |
|--------|------|----------|---------|---------|-------------|
| `PG_HOST` | string | ✗ | `localhost` | - | Database host |
| `PG_PORT` | int | ✗ | `5432` | - | Database port |
| `PG_USER` | string | ✓ | - | - | Database user |
| `PG_PASSWORD` | string | ✓ | - | - | Database password |
| `PG_DBNAME` | string | ✓ | - | - | Database name |

## Redis

Redis cache settings

| Name | Type | Required | Default | Example | Description |
|--------|------|----------|---------|---------|-------------|
| `REDIS_HOST` | string | ✗ | `localhost` | - | Redis host |
| `REDIS_PORT` | int | ✗ | `6379` | - | Redis port |
| `REDIS_PASSWORD` | string | ✗ | - | - | Redis password |
| `REDIS_DB` | int | ✗ | `0` | - | Redis database number | 