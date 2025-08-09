package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"secr-cli/internal/rules"
	"secr-cli/internal/scanner"
)

var staged bool

var rootCmd = &cobra.Command{
	Use:   "secr-cli [git-command] [args...]",
	Short: "Scan a git repo for secrets then run git command if clean",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cobraCmd *cobra.Command, args []string) {
		// Check if git is installed
		if _, err := exec.LookPath("git"); err != nil {
			fmt.Fprintln(os.Stderr, "Error: 'git' is not installed or not in PATH.")
			os.Exit(1)
		}

		// Check if inside a git repository
		gitCheck := exec.Command("git", "rev-parse", "--is-inside-work-tree")
		output, err := gitCheck.Output()
		if err != nil || strings.TrimSpace(string(output)) != "true" {
			fmt.Fprintln(os.Stderr, "Error: secr-cli must be run inside a Git repository.")
			os.Exit(1)
		}

		// Load rules and scan the repo
		ruleSet := rules.LoadRules()
		findings, err := scanner.ScanRepo(staged, ruleSet)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error during scan: %v\n", err)
			os.Exit(1)
		}

		if len(findings) > 0 {
			fmt.Println("Potential secrets detected:")
			for _, f := range findings {
				fmt.Printf(
					"  File: %s\n  Line: %d\n  Rule: %s\n  Content: %s\n\n",
					f.File, f.Line, f.RuleName, f.Content,
				)
			}

			if staged {
				fmt.Println("Secrets were found in staged (uncommitted) changes.")
				fmt.Println("To clean up:")
				for _, f := range findings {
					fmt.Printf("  git reset HEAD %s        # Unstage the affected file\n", f.File)
					fmt.Printf("  Edit '%s' to remove secrets\n", f.File)
					fmt.Printf("  git add %s               # Re-stage after fixing\n\n", f.File)
				}
			} else {
				fmt.Println("Secrets were found in the last commit.")
				fmt.Println("If not yet pushed:")
				for _, f := range findings {
					fmt.Println("  git reset --soft HEAD~1      # Undo the last commit")
					fmt.Printf("  Edit '%s' to remove secrets\n", f.File)
					fmt.Printf("  git add %s\n", f.File)
					fmt.Println("  git commit                   # Re-commit without secrets\n")
					break
				}
				fmt.Println("If already pushed:")
				for _, f := range findings {
					fmt.Printf("  git filter-repo --path %s --invert-paths\n", f.File)
					fmt.Println("  Then rotate any leaked credentials immediately.\n")
					break
				}
			}
			os.Exit(2)
		}

		gitCmd := exec.Command("git", args...)
		gitCmd.Stdin = os.Stdin
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr

		if err := gitCmd.Run(); err != nil {
			// Return git command's exit code if possible
			if exitError, ok := err.(*exec.ExitError); ok {
				os.Exit(exitError.ExitCode())
			}
			fmt.Fprintf(os.Stderr, "Error running git command: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	rootCmd.Flags().BoolVar(&staged, "staged", false, "Scan staged changes instead of HEAD")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
