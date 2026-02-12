# Secret Detection Rules

## Overview
This document lists all patterns detected by secr-cli, organized by category. Each rule has a severity level:
- 🔴 **HIGH** — Confirmed credential or key that should never be committed
- 🟡 **MEDIUM** — Likely sensitive, but may have legitimate uses
- 🔵 **LOW** — Potentially sensitive, higher chance of false positives

Run `secr-cli rules` to see all rules in your terminal.

---

## Cloud Credentials

| Rule | Severity | Pattern | Example |
|------|----------|---------|---------|
| AWS Access Key | 🔴 HIGH | `AKIA[0-9A-Z]{16}` | `AKIAIOSFODNN7EXAMPLE` |
| AWS Secret Key | 🔴 HIGH | `(?i)aws...['"][0-9a-zA-Z/+]{40}['"]` | `aws_secret_access_key = "wJalrX..."` |
| Google API Key | 🔴 HIGH | `AIza[0-9A-Za-z\-_]{35}` | `AIzaSyD_7fR2oX34X56Xy...` |
| Google OAuth Token | 🔴 HIGH | `ya29\.[a-zA-Z0-9\-_]+` | `ya29.a0ARrdaM...` |
| Azure Storage Key | 🔴 HIGH | `DefaultEndpointsProtocol=https;AccountName=...` | Full connection string |
| Heroku API Key | 🔴 HIGH | `[hH][eE][rR][oO][kK][uU]...UUID` | `heroku_api_key = "12345678-..."` |
| DigitalOcean Token | 🔴 HIGH | `dop_v1_[a-f0-9]{64}` | `dop_v1_abc123...` |
| Cloudflare API Key | 🔴 HIGH | `(?i)cloudflare...['"][a-z0-9]{37}['"]` | `cloudflare_key = "abc..."` |

## API Tokens

| Rule | Severity | Pattern | Example |
|------|----------|---------|---------|
| GitHub Token | 🔴 HIGH | `(ghp\|gho\|ghu\|ghs\|ghr)_[a-zA-Z0-9]{36}` | `ghp_3f6e6d9a1b2c...` |
| Slack Token | 🔴 HIGH | `xox[baprs]-[0-9a-zA-Z]{10,48}` | `xoxb-123456789012-...` |
| Stripe Secret Key | 🔴 HIGH | `sk_live_[0-9a-zA-Z]{24,99}` | `sk_live_abc123...` |
| Stripe Publishable Key | 🔵 LOW | `pk_live_[0-9a-zA-Z]{24,99}` | `pk_live_abc123...` |
| SendGrid API Key | 🔴 HIGH | `SG\.[a-zA-Z0-9\-_]{22,}\.…` | `SG.abc123.def456` |
| Twilio API Key | 🔴 HIGH | `SK[a-f0-9]{32}` | `SKabc123def456...` |
| npm Access Token | 🔴 HIGH | `npm_[a-zA-Z0-9]{36}` | `npm_abc123def456...` |
| PyPI API Token | 🔴 HIGH | `pypi-AgEIcHlwaS5vcmc...` | `pypi-AgEIcHlwaS5vcmc...` |
| Discord Bot Token | 🔴 HIGH | `[MN][A-Za-z\d]{23,}\.…` | `MTIzNDU2Nzg5MDEy...` |
| Telegram Bot Token | 🟡 MEDIUM | `[0-9]{8,10}:[a-zA-Z0-9_-]{35}` | `123456789:AAGB...` |
| Mailgun API Key | 🔴 HIGH | `key-[a-zA-Z0-9]{32}` | `key-abc123def456...` |
| Datadog API Key | 🟡 MEDIUM | `(?i)datadog...['"][a-f0-9]{32}['"]` | `datadog_key = "abc..."` |
| Shopify Token | 🔴 HIGH | `shpat_[a-fA-F0-9]{32}` | `shpat_abc123...` |
| Linear API Key | 🟡 MEDIUM | `lin_api_[a-zA-Z0-9]{40}` | `lin_api_abc123...` |
| OpenAI API Key | 🔴 HIGH | `sk-..T3BlbkFJ...` | `sk-abc123T3BlbkFJdef456` |
| Anthropic API Key | 🔴 HIGH | `sk-ant-api03-...` | `sk-ant-api03-abc123...` |
| Facebook Access Token | 🔴 HIGH | `EAACEdEose0cBA[0-9A-Za-z]+` | `EAACEdEose0cBAABC...` |
| Twitter API Key | 🟡 MEDIUM | `(?i)twitter...['"][0-9a-z]{35,44}['"]` | `twitter_key = "abc..."` |

## Cryptographic Material

| Rule | Severity | Pattern |
|------|----------|---------|
| RSA Private Key | 🔴 HIGH | `-----BEGIN RSA PRIVATE KEY-----` |
| EC Private Key | 🔴 HIGH | `-----BEGIN EC PRIVATE KEY-----` |
| DSA Private Key | 🔴 HIGH | `-----BEGIN DSA PRIVATE KEY-----` |
| SSH Private Key | 🔴 HIGH | `-----BEGIN OPENSSH PRIVATE KEY-----` |
| PGP Private Key | 🔴 HIGH | `-----BEGIN PGP PRIVATE KEY BLOCK-----` |
| Generic Private Key | 🔴 HIGH | `-----BEGIN PRIVATE KEY-----` |

## Database & Connection Strings

| Rule | Severity | Pattern | Example |
|------|----------|---------|---------|
| Database Connection String | 🔴 HIGH | `(jdbc:\|mongodb://\|...)...@host` | `postgres://user:pass@host` |
| Password in URL | 🔴 HIGH | `protocol://user:pass@host` | `https://admin:secret@server.com` |

## Authentication

| Rule | Severity | Pattern |
|------|----------|---------|
| JWT Token | 🟡 MEDIUM | `eyJ[base64].[base64].[base64]` |
| Basic Auth Credentials | 🟡 MEDIUM | `(?i)basic [base64]{5,100}` |
| Docker Registry Auth | 🔴 HIGH | `"auth"\s*:\s*"[base64]"` |

## Generic Patterns

| Rule | Severity | Notes |
|------|----------|-------|
| Generic API Key | 🟡 MEDIUM | Matches `api_key=`, `secret_key=`, `access_token=`, etc. |
| Env File Secret | 🟡 MEDIUM | Matches `PASSWORD=`, `SECRET=`, `TOKEN=`, etc. in env-style files |

## Payment Information

| Rule | Severity | Notes |
|------|----------|-------|
| Credit Card Number | 🔴 HIGH | Visa, Mastercard, Amex, Diners, Discover, JCB |

---

## Adding New Rules

1. Add pattern to `internal/rules/rules.go` with a `Severity` level
2. Document here with pattern, example, and severity
3. Run `secr-cli rules` to verify it shows up
4. Submit PR for review

> **Note:** Some patterns are intentionally broad to catch variants while minimizing false negatives. Use `--severity HIGH` to filter noise.