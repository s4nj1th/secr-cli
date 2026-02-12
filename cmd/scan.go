package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"secr-cli/internal/rules"
	"secr-cli/internal/scanner"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	scanShow        bool
	scanStagedOnly  bool
	scanJSON        bool
	scanNoGitignore bool
	scanWorkers     int
	scanSeverity    string
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan the repository for secrets",
	Long:  "Scan staged, unstaged, and working directory files for secrets like API keys, tokens, and credentials.",
	Run: func(cmd *cobra.Command, args []string) {
		runScan(scanShow, scanStagedOnly, scanJSON, scanNoGitignore, scanWorkers, scanSeverity)
	},
}

func init() {
	scanCmd.Flags().BoolVarP(&scanShow, "show", "s", false, "Display secret content in output")
	scanCmd.Flags().BoolVar(&scanStagedOnly, "staged-only", false, "Only scan staged changes")
	scanCmd.Flags().BoolVar(&scanJSON, "json", false, "Output findings as JSON")
	scanCmd.Flags().BoolVar(&scanNoGitignore, "no-gitignore", false, "Scan all files, ignoring .gitignore rules")
	scanCmd.Flags().IntVar(&scanWorkers, "workers", 0, "Number of concurrent workers (default: number of CPUs)")
	scanCmd.Flags().StringVar(&scanSeverity, "severity", "", "Minimum severity to report (LOW, MEDIUM, HIGH)")
}

func runScan(showContent, stagedOnly, jsonOutput, noGitignore bool, workers int, minSeverity string) {
	ruleSet := rules.LoadRules()

	opts := scanner.ScanOptions{
		StagedOnly:  stagedOnly,
		NoGitignore: noGitignore,
		Workers:     workers,
	}

	findings, err := scanner.ScanRepo(ruleSet, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during scan: %v\n", err)
		os.Exit(1)
	}

	if minSeverity != "" {
		findings = filterBySeverity(findings, strings.ToUpper(minSeverity))
	}

	if jsonOutput {
		printJSON(findings, showContent)
		if len(findings) > 0 {
			os.Exit(2)
		}
		return
	}

	if len(findings) > 0 {
		printFindings(findings, showContent)
		os.Exit(2)
	} else {
		fmt.Println(clearBg(" No secrets detected! "))
	}
}

func filterBySeverity(findings []scanner.Finding, minSev string) []scanner.Finding {
	levels := map[string]int{"LOW": 1, "MEDIUM": 2, "HIGH": 3}
	minLevel, ok := levels[minSev]
	if !ok {
		return findings
	}

	var filtered []scanner.Finding
	for _, f := range findings {
		if level, ok := levels[f.Severity]; ok && level >= minLevel {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

func printFindings(findings []scanner.Finding, showContent bool) {
	fmt.Println(alertBg(" Potential secrets detected! "))

	severityColor := map[string]func(a ...interface{}) string{
		"HIGH":   color.New(color.FgRed, color.Bold).SprintFunc(),
		"MEDIUM": color.New(color.FgYellow).SprintFunc(),
		"LOW":    color.New(color.FgBlue).SprintFunc(),
	}

	grouped := make(map[string][]scanner.Finding)
	for _, f := range findings {
		grouped[f.Type] = append(grouped[f.Type], f)
	}

	for _, t := range []string{"staged", "unstaged", "working"} {
		if finds, ok := grouped[t]; ok {
			header := color.New(color.Bold, color.Underline).Sprintf("\n%s changes:", strings.ToUpper(t[:1])+t[1:])
			fmt.Println(header)
			for _, f := range finds {
				sevFn, ok := severityColor[f.Severity]
				if !ok {
					sevFn = fmt.Sprint
				}
				fmt.Printf("  File: %s\n  Line: %d\n  Rule: %s [%s]\n",
					f.File, f.Line, f.RuleName, sevFn(f.Severity),
				)
				if showContent {
					fmt.Printf("  Content: %s\n\n", f.Content)
				} else {
					fmt.Println()
				}
			}
		}
	}
}

func printJSON(findings []scanner.Finding, showContent bool) {
	type jsonFinding struct {
		File     string `json:"file"`
		Line     int    `json:"line"`
		Rule     string `json:"rule"`
		Severity string `json:"severity"`
		Type     string `json:"type"`
		Content  string `json:"content,omitempty"`
	}

	var output []jsonFinding
	for _, f := range findings {
		jf := jsonFinding{
			File:     f.File,
			Line:     f.Line,
			Rule:     f.RuleName,
			Severity: f.Severity,
			Type:     f.Type,
		}
		if showContent {
			jf.Content = f.Content
		}
		output = append(output, jf)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(output)
}
