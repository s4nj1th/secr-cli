package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"secr-cli/internal/rules"
	"secr-cli/internal/scanner"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var dangerousCommands = map[string]bool{
	"commit":      true,
	"push":        true,
	"merge":       true,
	"rebase":      true,
	"stash":       true,
	"am":          true,
	"cherry-pick": true,
}

var gitCmd = &cobra.Command{
	Use:   "git [git-args...]",
	Short: "Drop-in git replacement with automatic secret scanning",
	Long: `Use secr-cli as a drop-in replacement for git.

Safe commands (status, log, diff, branch, etc.) pass through instantly.
Dangerous commands (commit, push, merge, etc.) trigger a secret scan first.

Setup:
  secr-cli wrap        # adds alias to your shell config
  secr-cli unwrap      # removes the alias

Manual setup:
  alias git='secr-cli git'

After this, use git normally:
  git status           # instant passthrough
  git commit -m "msg"  # scans for secrets first
  git push origin main # scans for secrets first
  git log              # instant passthrough`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {

			runGit(args)
			return
		}

		subcommand := findGitSubcommand(args)

		if dangerousCommands[subcommand] {

			dim := color.New(color.Faint).SprintFunc()
			fmt.Fprintf(os.Stderr, "%s\n", dim("secr-cli: scanning before "+subcommand+"..."))

			ruleSet := rules.LoadRules()
			opts := scanner.ScanOptions{}

			if subcommand == "commit" {
				opts.StagedOnly = true
			}

			findings, err := scanner.ScanRepo(ruleSet, opts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "secr-cli: scan error: %v\n", err)
				os.Exit(1)
			}

			if len(findings) > 0 {
				printFindings(findings, false)
				fmt.Fprintln(os.Stderr)
				fmt.Fprintln(os.Stderr, alertBg(" Secrets detected! ")+" git "+subcommand+" aborted.")
				fmt.Fprintln(os.Stderr, "Fix the issues above, then try again.")
				if subcommand == "commit" {
					fmt.Fprintln(os.Stderr, "To bypass: git commit --no-verify (NOT recommended)")
				}
				os.Exit(2)
			}

			fmt.Fprintf(os.Stderr, "%s\n", dim("secr-cli: clean ✓"))
		}

		runGit(args)
	},
}

func findGitSubcommand(args []string) string {

	flagsWithArg := map[string]bool{
		"-C": true, "-c": true, "--git-dir": true,
		"--work-tree": true, "--namespace": true,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			if flagsWithArg[arg] && i+1 < len(args) {
				i++
			}
			continue
		}
		return arg
	}
	return ""
}

func runGit(args []string) {
	gitBin, err := findRealGit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "secr-cli: %v\n", err)
		os.Exit(1)
	}

	gitExec := exec.Command(gitBin, args...)
	gitExec.Stdin = os.Stdin
	gitExec.Stdout = os.Stdout
	gitExec.Stderr = os.Stderr

	if err := gitExec.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}

func findRealGit() (string, error) {
	self, _ := os.Executable()
	selfBase := filepath.Base(self)

	commonPaths := []string{
		"/usr/bin/git",
		"/usr/local/bin/git",
		"/opt/homebrew/bin/git",
	}
	for _, p := range commonPaths {
		if info, err := os.Stat(p); err == nil && !info.IsDir() {

			realPath, _ := filepath.EvalSymlinks(p)
			if filepath.Base(realPath) != selfBase {
				return p, nil
			}
		}
	}

	pathEnv := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(pathEnv) {
		candidate := filepath.Join(dir, "git")
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			realPath, _ := filepath.EvalSymlinks(candidate)
			realSelf, _ := filepath.EvalSymlinks(self)
			if realPath != realSelf {
				return candidate, nil
			}
		}
	}

	return "", fmt.Errorf("could not find git binary (is git installed?)")
}
