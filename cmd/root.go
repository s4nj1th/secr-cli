package cmd

import (
	"fmt"
	"os"

	"secr-cli/internal/rules"
	"secr-cli/internal/scanner"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var show bool

var (
	alertBg = color.New(color.BgRed, color.FgBlack).SprintFunc()
	clearBg = color.New(color.BgGreen, color.FgBlack).SprintFunc()
)

var rootCmd = &cobra.Command{
	Use:   "secr-cli",
	Short: "A lightning-fast secret scanner for Git repositories",
	Long: `secr-cli scans Git repositories for secrets like API keys, tokens, 
private keys, and credentials. It supports staged, unstaged, and 
working directory scans with concurrent processing.

Quick start:
  secr-cli                          Scan the repo for secrets
  secr-cli scan --staged-only       Scan only staged changes
  secr-cli scan --json              Output findings as JSON
  secr-cli hook install             Install pre-commit hook
  secr-cli git commit -m "msg"      Scan before running git commit`,
	Run: func(cobraCmd *cobra.Command, args []string) {
		runScan(show, false, false, false, 0, "")
	},
}

var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of secr-cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("secr-cli %s\n", version)
	},
}

var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "List all detection rules",
	Run: func(cmd *cobra.Command, args []string) {
		ruleSet := rules.LoadRules()

		severityColor := map[string]func(a ...interface{}) string{
			"HIGH":   color.New(color.FgRed, color.Bold).SprintFunc(),
			"MEDIUM": color.New(color.FgYellow).SprintFunc(),
			"LOW":    color.New(color.FgBlue).SprintFunc(),
		}

		fmt.Printf("%-35s %-8s  %s\n", "RULE", "SEVERITY", "PATTERN")
		fmt.Println(color.New(color.Faint).Sprint("─────────────────────────────────────────────────────────────────────────────"))
		for _, r := range ruleSet {
			sevFn, ok := severityColor[string(r.Severity)]
			if !ok {
				sevFn = fmt.Sprint
			}
			pattern := r.Pattern.String()
			if len(pattern) > 40 {
				pattern = pattern[:37] + "..."
			}
			fmt.Printf("%-35s %-8s  %s\n", r.Name, sevFn(string(r.Severity)), pattern)
		}
		fmt.Printf("\nTotal: %d rules\n", len(ruleSet))
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show scan status and repository info",
	Run: func(cmd *cobra.Command, args []string) {
		ruleSet := rules.LoadRules()

		opts := scanner.ScanOptions{}
		findings, err := scanner.ScanRepo(ruleSet, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Rules loaded:       %d\n", len(ruleSet))

		staged, unstaged, working := 0, 0, 0
		for _, f := range findings {
			switch f.Type {
			case "staged":
				staged++
			case "unstaged":
				unstaged++
			case "working":
				working++
			}
		}

		fmt.Printf("Findings (staged):  %d\n", staged)
		fmt.Printf("Findings (unstaged):%d\n", unstaged)
		fmt.Printf("Findings (working): %d\n", working)
		fmt.Printf("Total findings:     %d\n", len(findings))
	},
}

func Execute() {
	rootCmd.Flags().BoolVarP(&show, "show", "s", false, "Display secret content in output")

	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(hookCmd)
	rootCmd.AddCommand(gitCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(rulesCmd)
	rootCmd.AddCommand(statusCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
