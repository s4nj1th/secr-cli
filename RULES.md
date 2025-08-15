# Secret Detection Rules

## Overview
This document lists all patterns detected by secr-cli, organized by category. Each rule has:
- A descriptive name
- Example matches
- False positive guidance
- Severity level

## Cloud Credentials

### AWS Access Key
- **Pattern:** `AKIA[0-9A-Z]{16}`
- **Example:** `AKIAIOSFODNN7EXAMPLE`
- **Notes:** AWS root account credentials

### AWS Secret Key
- **Pattern:** `(?i)aws(.{0,20})?(secret)?(.{0,20})?['"][0-9a-zA-Z/+]{40}['"]`
- **Example:** `aws_secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"`

### Google API Key
- **Pattern:** `AIza[0-9A-Za-z\-_]{35}`
- **Example:** `AIzaSyD_7fR2oX34X56XyEXAMPLEKEYq1o`

## API Tokens

### GitHub Token
- **Pattern:** `(ghp|gho|ghu|ghs|ghr)_[a-zA-Z0-9]{36}`
- **Example:** `ghp_3f6e6d9a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6`

### Slack Token
- **Pattern:** `xox[baprs]-[0-9a-zA-Z]{10,48}`
- **Example:** `xoxb-123456789012-123456789012-abcDEfGhIjkLmNoPQRsTuVwXyZ`

## Cryptographic Material

### RSA Private Key
- **Pattern:** `-----BEGIN RSA PRIVATE KEY-----`
- **Notes:** Looks for PEM headers

### SSH Private Key
- **Pattern:** `-----BEGIN OPENSSH PRIVATE KEY-----`

## Database Credentials

### PostgreSQL Connection String
- **Pattern:** `postgres(ql)?://[^:]+:[^@]+@`
- **Example:** `postgres://user:password@localhost:5432/db`

### MongoDB Connection String
- **Pattern:** `mongodb(\\+srv)?://[^:]+:[^@]+@`
- **Example:** `mongodb://admin:password@cluster0.example.com:27017`

## Payment Information

### Credit Card Numbers
- **Pattern:** `\b(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|3[47][0-9]{13})\b`
- **False Positives:** Test card numbers (like 4242...), invoice numbers

## Adding New Rules
1. Add pattern to `internal/rules/rules.go`
2. Document here with:
   - Exact pattern
   - Example match
   - Severity justification
3. Submit PR for review

> **Note:** Some patterns are intentionally vague to catch variants while minimizing false negatives.