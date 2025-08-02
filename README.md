<div align="center">
<h1>secr-cli</h1>
</div>

**`secr-cli`** is a fast and minimal command-line tool written in Go for scanning Git repositories for sensitive information such as API keys, tokens, and private keys. It works on both committed (`HEAD`) and staged changes.

## Features

* Scans either:

  * The latest commit (`HEAD`)
  * Staged (uncommitted) changes
* Detects:

  * AWS credentials
  * Generic API keys
  * Private keys
  * Slack tokens
* Outputs filename, line number, and rule matched
* No external dependencies beyond Go standard library

## Installation

### Build from Source

**Requirements**: Go 1.21+

Clone the repository:

```bash
git clone https://github.com/s4nj1th/secr-cli.git
cd secr-cli
```

Build the CLI:

```bash
go build -o secr-cli ./main.go
```

This will create a `secr-cli` binary in the current directory.

To verify the build:

```bash
./secr-cli --help
```

## Usage

### Scan latest commit (HEAD)

```bash
./secr-cli
```

### Scan staged (uncommitted) changes

```bash
./secr-cli --staged
```

## Example Output

```
[Generic API Key] config/dev.env:12: SECRET_KEY=abcdef1234567890abcdef1234567890
[AWS Access Key] main.go:89: AKIAIOSFODNN7EXAMPLE
```

## Patterns Detected

* AWS Access Key (e.g. `AKIA...`)
* AWS Secret Key
* RSA/DSA/EC Private Keys
* Generic API Keys
* Slack Tokens

Pattern definitions can be found in [`internal/rules/rules.go`](./internal/rules/rules.go)

## Contributing

Open issues or submit pull requests to:

* Add more detection rules
* Improve CLI usability
* Add `.gitignore` support
* Support JSON or SARIF output

## License

MIT. See [LICENSE](./LICENSE).

Let me know if you want badges (Go version, license, etc.), Docker usage instructions, or packaging steps (e.g., `make`, release automation).
