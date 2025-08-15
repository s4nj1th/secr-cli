<div align="center">
<h1>secr-cli</h1>
<p>A lightning-fast secret scanner for Git repositories</p>
</div>

**`secr-cli`** is a fast and minimal command-line tool written in Go for scanning Git repositories for sensitive information such as API keys, tokens, and private keys. It works on both committed (`HEAD`), staged changes, unstaged changes, and can also act as a wrapper around Git commands to enforce secret scanning before any Git operation.

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

1. Clone the repository
2. Build the binary
3. Move the binary to your `PATH`

```bash
git clone https://github.com/s4nj1th/secr-cli
cd secr-cli
sudo make install
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
# Find where the secrets are
secr-cli

# Show secret content (careful!)
secr-cli --show
```

### Using secr-cli as a Git wrapper

You can configure `secr-cli` to automatically scan for secrets before running any Git command. To do this, create a shell alias:

```bash
alias git='secr-cli && git'
```

Add this line to your shell configuration file (`~/.bashrc`, `~/.zshrc`, etc.) to make it persistent.

After this, when you run any `git` command, for example:

```bash
git commit -m "fix typo"
```

`secr-cli` will first scan for secrets in the repository (HEAD or staged, depending on flags), and if none are found, it will forward the command to the actual Git binary.

If secrets are detected, the Git command is aborted and you will be shown details about the secrets found along with remediation steps.

## Patterns Detected

See all patterns in [RULES](RULES.md).

## Contributing

Open issues or submit pull requests to:

* Add more detection rules
* Improve CLI usability
* Add `.gitignore` support
* Support JSON or SARIF output

## License

This project is licensed under the GNU General Public License v3.0.
See the [COPYING](./COPYING) file for details.

## Contributing
We welcome contributions! Please see [CONTRIBUTING](CONTRIBUTING.md).

