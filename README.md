# KrakenD No Redirects Plugin

[![CI](https://github.com/ulbwa/krakend-no-redirects-plugin/actions/workflows/ci.yml/badge.svg)](https://github.com/ulbwa/krakend-no-redirects-plugin/actions/workflows/ci.yml)
[![Release](https://github.com/ulbwa/krakend-no-redirects-plugin/actions/workflows/cd.yml/badge.svg)](https://github.com/ulbwa/krakend-no-redirects-plugin/actions/workflows/cd.yml)
[![Go Version](https://img.shields.io/badge/Go-1.25.3-blue.svg)](https://golang.org)

A KrakenD HTTP client plugin that disables automatic redirect following. This plugin forces the HTTP client to return redirect responses (3xx) directly to the caller instead of automatically following them.

## Features

- Disables automatic HTTP redirect following
- Returns 3xx responses with `Location` headers intact
- Preserves all response headers and status codes
- Compatible with KrakenD 2.x
- Available for Linux AMD64 and ARM64 architectures

## Compatibility

- **KrakenD:** 2.x
- **Go:** 1.25.3
- **Platforms:** Linux (amd64, arm64)

## Installation

### Download Pre-built Binary

Download the latest release for your architecture:

```bash
# For AMD64
wget https://github.com/ulbwa/krakend-no-redirects-plugin/releases/latest/download/no-redirects-linux-amd64.so

# For ARM64
wget https://github.com/ulbwa/krakend-no-redirects-plugin/releases/latest/download/no-redirects-linux-arm64.so
```

### Build from Source

**Note:** The plugin must be built on Linux as it requires CGO and Linux-specific compilation. 

For local development and testing, use the provided Docker setup:

```bash
git clone https://github.com/ulbwa/krakend-no-redirects-plugin.git
cd krakend-no-redirects-plugin
make build
```

The compiled plugin will be available in the `bin/` directory.

## Usage

### Docker

Mount the plugin file to `/etc/krakend/plugins/` in your KrakenD container:

```bash
docker run -p 8080:8080 \
  -v $PWD/no-redirects-linux-amd64.so:/etc/krakend/plugins/no-redirects.so:ro \
  -v $PWD/krakend.json:/etc/krakend/krakend.json:ro \
  krakend:latest run -c /etc/krakend/krakend.json
```

### Docker Compose

```yaml
version: '3'
services:
  krakend:
    image: krakend:latest
    ports:
      - "8080:8080"
    volumes:
      - ./no-redirects-linux-amd64.so:/etc/krakend/plugins/no-redirects.so:ro
      - ./krakend.json:/etc/krakend/krakend.json:ro
    command: ["run", "-c", "/etc/krakend/krakend.json"]
```

### Self-Hosted Installation

1. Copy the plugin to your KrakenD plugins directory:

```bash
sudo cp no-redirects-linux-amd64.so /usr/local/lib/krakend/plugins/no-redirects.so
sudo chmod 644 /usr/local/lib/krakend/plugins/no-redirects.so
```

2. Update your KrakenD configuration to load the plugin.

## Configuration

### Enable Plugin Globally

Add the plugin to your `krakend.json` configuration:

```json
{
  "version": 3,
  "plugin": {
    "folder": "/etc/krakend/plugins/",
    "pattern": ".so"
  }
}
```

### Use Plugin for Specific Endpoints

Configure the plugin as an HTTP client for endpoints that require redirect handling:

```json
{
  "version": 3,
  "endpoints": [
    {
      "endpoint": "/redirect-test",
      "method": "GET",
      "output_encoding": "no-op",
      "backend": [
        {
          "url_pattern": "/redirect-to?url=https%3A%2F%2Fgoogle.com&status_code=302",
          "host": ["https://httpbin.org"],
          "encoding": "no-op",
          "extra_config": {
            "plugin/http-client": {
              "name": "krakend-no-redirects"
            }
          }
        }
      ]
    }
  ]
}
```

<details>
<summary>Complete Example</summary>

```json
{
  "version": 3,
  "plugin": {
    "folder": "/etc/krakend/plugins/",
    "pattern": ".so"
  },
  "endpoints": [
    {
      "endpoint": "/api/redirect",
      "method": "GET",
      "output_encoding": "no-op",
      "backend": [
        {
          "url_pattern": "/redirect-to?url=https%3A%2F%2Fgoogle.com&status_code=302",
          "host": ["https://httpbin.org"],
          "encoding": "no-op",
          "extra_config": {
            "plugin/http-client": {
              "name": "krakend-no-redirects"
            }
          }
        }
      ]
    }
  ]
}
```

When you call `GET /api/redirect`, the plugin will return the 302/301 response directly instead of following the redirect.

</details>

## How It Works

The plugin replaces KrakenD's default HTTP client with a custom client that sets `CheckRedirect` to return `http.ErrUseLastResponse`. This prevents the Go HTTP client from automatically following redirects and returns the redirect response to the caller.

## Development

### Prerequisites

- Go 1.25.3 or later
- Make
- GCC (for CGO)

### Build

```bash
make build
```

The compiled plugin will be available in the `bin/` directory.

**Note:** This plugin does not support cross-compilation. You must build on the target platform (AMD64 or ARM64) to generate the corresponding `.so` file. The build process automatically detects the host architecture and produces a binary for that platform only.

### Clean Build Artifacts

```bash
make clean
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## References

- [Выпускайте Кракена: опыт использования KrakenD](https://habr.com/ru/companies/ru_mts/articles/716512/)
- [KrakenD Plugin Documentation](https://www.krakend.io/docs/extending/http-client-plugins/)
- [KrakenD Configuration](https://www.krakend.io/docs/configuration/structure/)
