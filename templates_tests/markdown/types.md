# Environment Variables Documentation

## App

Application settings

| Name | Type | Required | Default | Example | Description |
|--------|------|----------|---------|---------|-------------|
| `ENV` | string | ✓ | - | `development` | Application environment (Possible values: development, staging, production) |
| `API_URL` | *url.URL | ✓ | - | `https://api.example.com` | API endpoint |
| `REQUEST_TIMEOUT` | time.Duration | ✗ | `30s` | - | API request timeout |
| `RESPONSE_TIMEOUT` | time.Duration | ✗ | `30s` | - | API response timeout |
| `ALLOWED_IPS` | []net.IP | ✗ | - | `127.0.0.1,192.168.1.1` | List of allowed IP addresses |

## Custom Types

| Type | Import Path | Description | Possible Values |
|----|------------|-------------|----------------|
| `Environment` | - | Application environment | `development`, `staging`, `production` |
| `URL` | `net/url` | URL type | - |
| `Duration` | `time` | Duration type | - |
| `IpAddresses` | `net` | IP addresses type | - | 