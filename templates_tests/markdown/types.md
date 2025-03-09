# Environment Variables Documentation

## App

Application settings

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `ENV` | string | ✓ | - | Application environment (Possible values: development, staging, production) |
| `API_URL` | *url.URL | ✓ | - | API endpoint |
| `REQUEST_TIMEOUT` | time.Duration | ✗ | `30s` | API request timeout |
| `RESPONSE_TIMEOUT` | time.Duration | ✗ | `30s` | API response timeout |
| `ALLOWED_IPS` | []net.IP | ✗ | - | List of allowed IP addresses |

## Types

### Environment

Application environment

Possible values:
- `development`
- `staging`
- `production`

### URL

URL type

### Duration

Duration type

### IpAddresses

IP addresses type 