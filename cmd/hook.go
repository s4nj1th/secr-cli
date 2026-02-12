package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

const preCommitHook = `#!/bin/sh
# secr-cli pre-commit hook
# Automatically installed by: secr-cli hook install

echo "secr-cli: scanning for secrets before commit..."
secr-cli scan --staged-only
STATUS=$?
if [ $STATUS -ne 0 ]; then
    echo ""
    echo "secr-cli: secrets detected! Commit aborted."
    echo "secr-cli: fix the issues above, then try again."
    echo "secr-cli: to bypass this check (NOT recommended), use: git commit --no-verify"
    exit 1
fi
`

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Manage Git pre-commit hooks",
	Long:  "Install or uninstall a Git pre-commit hook that automatically scans for secrets before each commit.",
}

var hookInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the pre-commit hook",
	Run: func(cmd *cobra.Command, args []string) {
		hookPath, err := getHookPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if _, err := os.Stat(hookPath); err == nil {
			fmt.Fprintf(os.Stderr, "Warning: pre-commit hook already exists at %s\n", hookPath)
			fmt.Fprintf(os.Stderr, "Use --force to overwrite.\n")

			force, _ := cmd.Flags().GetBool("force")
			if !force {
				os.Exit(1)
			}
		}

		if err := os.MkdirAll(filepath.Dir(hookPath), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating hooks directory: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(hookPath, []byte(preCommitHook), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing hook: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(clearBg(" Pre-commit hook installed! "))
		fmt.Printf("Hook path: %s\n", hookPath)
		fmt.Println("Secrets will be scanned automatically before each commit.")
	},
}

var hookUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove the pre-commit hook",
	Run: func(cmd *cobra.Command, args []string) {
		hookPath, err := getHookPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if _, err := os.Stat(hookPath); os.IsNotExist(err) {
			fmt.Println("No pre-commit hook found. Nothing to remove.")
			return
		}

		content, err := os.ReadFile(hookPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading hook: %v\n", err)
			os.Exit(1)
		}

		if string(content) != preCommitHook {
			fmt.Fprintf(os.Stderr, "Warning: pre-commit hook was not installed by secr-cli.\n")
			fmt.Fprintf(os.Stderr, "Use --force to remove it anyway.\n")

			force, _ := cmd.Flags().GetBool("force")
			if !force {
				os.Exit(1)
			}
		}

		if err := os.Remove(hookPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing hook: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(clearBg(" Pre-commit hook removed! "))
	},
}

func init() {
	hookInstallCmd.Flags().Bool("force", false, "Overwrite existing pre-commit hook")
	hookUninstallCmd.Flags().Bool("force", false, "Remove hook even if not installed by secr-cli")

	hookCmd.AddCommand(hookInstallCmd)
	hookCmd.AddCommand(hookUninstallCmd)
}

func getHookPath() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository")
	}

	gitDir := filepath.Clean(string(out[:len(out)-1]))
	return filepath.Join(gitDir, "hooks", "pre-commit"), nil
}
