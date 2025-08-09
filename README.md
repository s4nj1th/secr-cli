<div align="center">
<h1>secr-cli</h1>
</div>

**`secr-cli`** is a fast and minimal command-line tool written in Go for scanning Git repositories for sensitive information such as API keys, tokens, and private keys. It works on both committed (`HEAD`) and staged changes, and can also act as a wrapper around Git commands to enforce secret scanning before any Git operation.

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
* Can be used as a Git command wrapper to scan secrets automatically before running any Git command

## Installation

### Option 1: Install from Releases

Download the latest precompiled binary for your platform from the [Releases](https://github.com/s4nj1th/secr-cli/releases) page.

Make it executable and move it into your `$PATH`. For example:

```bash
chmod +x secr-cli
sudo mv secr-cli /usr/local/bin/
```

### Option 2: Build from Source

**Requirements:** Go 1.21+

Clone the repository:

```bash
make install
```

This will compile and copy `secr-cli` to `/usr/local/bin/` (you might need `sudo`).

Verify the build:

```bash
secr-cli --help
```

## Usage

### Basic secret scan commands

Scan the latest commit (`HEAD`):

```bash
secr-cli
```

Scan staged (uncommitted) changes:

```bash
secr-cli --staged
```

### Using secr-cli as a Git wrapper

You can configure `secr-cli` to automatically scan for secrets before running any Git command. To do this, create a shell alias:

```bash
alias git='secr-cli'
```

Add this line to your shell configuration file (`~/.bashrc`, `~/.zshrc`, etc.) to make it persistent.

After this, when you run any `git` command, for example:

```bash
git commit -m "fix typo"
```

`secr-cli` will first scan for secrets in the repository (HEAD or staged, depending on flags), and if none are found, it will forward the command to the actual Git binary.

If secrets are detected, the Git command is aborted and you will be shown details about the secrets found along with remediation steps.

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

This project is licensed under the GNU General Public License v3.0.
See the [COPYING](./COPYING) file for details.
