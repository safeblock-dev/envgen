# Environment Variables Documentation

## App

Application settings

| Name | Type | Required | Default | Example | Description |
|--------|------|----------|---------|---------|-------------|
| `ENV` | [`Environment`](#custom-types) (string) | ✓ | - | `development` | Application environment (Possible values: development, staging, production) |
| `API_URL` | [`URL`](#custom-types) (*url.URL) | ✓ | - | `https://api.example.com` | API endpoint |
| `REQUEST_TIMEOUT` | [`Duration`](#custom-types) (time.Duration) | ✗ | `30s` | - | API request timeout |
| `RESPONSE_TIMEOUT` | [`Duration`](#custom-types) (time.Duration) | ✗ | `30s` | - | API response timeout |
| `ALLOWED_IPS` | [`IpAddresses`](#custom-types) ([]net.IP) | ✗ | - | `127.0.0.1,192.168.1.1` | List of allowed IP addresses |

## Custom Types

| Name | Type | Import Path | Description | Possible Values |
|----|------|------------|-------------|----------------|
| `Environment` | string | - | Application environment | `development`, `staging`, `production` |
| `URL` | *url.URL | `net/url` | URL type | - |
| `Duration` | time.Duration | `time` | Duration type | - |
| `IpAddresses` | []net.IP | `net` | IP addresses type | - | 