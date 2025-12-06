# Preflight CLI

A command-line tool that scans your codebase for launch readiness. Identifies missing configuration, integration issues, security concerns, SEO metadata gaps, and other common mistakes before you deploy to production.

## Installation

### Homebrew (macOS/Linux)

```bash
brew install phillips-jon/preflight/preflight
```

### Shell Script

```bash
curl -sSL https://raw.githubusercontent.com/phillips-jon/preflight.sh/main/install.sh | sh
```

### Manual Download

Download the latest release from [GitHub Releases](https://github.com/phillips-jon/preflight.sh/releases).

## Quick Start

```bash
# Initialize in your project directory
cd your-project
preflight init

# Run all checks
preflight scan

# Run in CI mode with JSON output
preflight scan --ci --format json
```

## What It Checks

| Check | Description |
|-------|-------------|
| **ENV Parity** | Compares `.env` and `.env.example` for missing variables |
| **Health Endpoint** | Verifies `/health` is reachable on staging/production |
| **Stripe Webhook** | Confirms webhook endpoint is accessible |
| **Sentry Init** | Detects Sentry initialization in codebase |
| **Plausible Analytics** | Finds Plausible script tag in templates |
| **SEO Metadata** | Checks for title, description, and Open Graph tags |
| **Security Headers** | Validates HSTS, CSP, and other security headers |
| **Secret Scanning** | Finds leaked API keys and credentials |

## Configuration

Preflight uses a `preflight.yml` file in your project root:

```yaml
projectName: my-app
stack: rails  # rails, next, node, laravel, static

urls:
  staging: "https://staging.example.com"
  production: "https://example.com"

services:
  stripe:
    declared: true
  sentry:
    declared: true

checks:
  envParity:
    enabled: true
    envFile: ".env"
    exampleFile: ".env.example"

  healthEndpoint:
    enabled: true
    path: "/health"

  stripeWebhook:
    enabled: true
    url: "https://api.example.com/webhooks/stripe"

  seoMeta:
    enabled: true
    mainLayout: "app/views/layouts/application.html.erb"
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All checks passed |
| 1 | Warnings only |
| 2 | Errors found |

## Supported Stacks

- Ruby on Rails
- Next.js
- Node.js (Express, etc.)
- Laravel
- Static sites

## CI Integration

```yaml
# GitHub Actions example
- name: Run Preflight
  run: |
    curl -sSL https://raw.githubusercontent.com/phillips-jon/preflight.sh/main/install.sh | sh
    preflight scan --ci --format json
```

## License

MIT
