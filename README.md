<div align="center">
<h1>secr-cli</h1>
<p>A lightning-fast secret scanner for Git repositories</p>
</div>

**`secr-cli`** is a fast and minimal command-line tool written in Go for scanning Git repositories for sensitive information such as API keys, tokens, and private keys. It uses **goroutine-based concurrent scanning** and respects `.gitignore` rules out of the box.

## Features

- **40+ detection rules** for API keys, tokens, credentials, and secrets
- **Severity levels** (HIGH / MEDIUM / LOW) for all rules
- **`.gitignore`-aware** — automatically skips ignored files
- **Concurrent scanning** — goroutine worker pool for fast file scanning
- **Git integration** — pre-commit hook management and git command passthrough
- **JSON output** — machine-readable output for CI/CD pipelines
- **Staged-only mode** — scan only what you're about to commit

## Installation

### Option 1: Install from Releases

Download the latest precompiled binary for your platform from the [Releases](https://github.com/s4nj1th/secr-cli/releases) page.

```bash
chmod +x secr-cli
sudo mv secr-cli /usr/local/bin/
```

### Option 2: Build from Source

**Requirements:** Go 1.21+

```bash
git clone https://github.com/s4nj1th/secr-cli
cd secr-cli
sudo make install
```

Verify:

```bash
secr-cli --help
```

## Usage

### Quick Scan

```bash
# Scan the repo (staged + unstaged + working directory)
secr-cli

# Show secret content (careful!)
secr-cli --show
```

### Scan Subcommand

```bash
# Scan only staged changes (great for pre-commit)
secr-cli scan --staged-only

# Output as JSON (for CI/CD)
secr-cli scan --json

# Filter by severity
secr-cli scan --severity HIGH

# Scan everything, ignore .gitignore rules
secr-cli scan --no-gitignore

# Control concurrency
secr-cli scan --workers 8
```

### Pre-Commit Hook

Install a Git pre-commit hook that automatically scans for secrets:

```bash
# Install the hook
secr-cli hook install

# Remove the hook
secr-cli hook uninstall
```

After installation, every `git commit` will automatically scan staged changes first.

### Git Passthrough

Run any Git command with an automatic secret scan:

```bash
secr-cli git commit -m "my changes"
secr-cli git push origin main
secr-cli git merge feature-branch
```

If secrets are detected, the Git command is aborted.

### Other Commands

```bash
# List all detection rules with severity
secr-cli rules

# Show scan status summary
secr-cli status

# Print version
secr-cli version
```

### Shell Alias (Optional)

You can also alias `git` to always scan first:

```bash
alias git='secr-cli git'
```

Add to your shell config (`~/.bashrc`, `~/.zshrc`) to make it persistent.

## Patterns Detected

See all patterns in [RULES](RULES.md), or run `secr-cli rules` to list them in terminal.

**Categories:** Cloud Credentials, API Tokens, Cryptographic Material, Database Credentials, Authentication, Payment Information, Generic Patterns.

## Contributing

We welcome contributions! Please see [CONTRIBUTING](CONTRIBUTING.md).

Open issues or submit pull requests to:

- Add more detection rules
- Improve CLI usability
- Support SARIF output
- Add custom rule configuration

## License

This project is licensed under the GNU General Public License v3.0.
See the [COPYING](./COPYING) file for details.
