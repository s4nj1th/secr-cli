package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/fatih/color"
	"secr-cli/internal/rules"
	"secr-cli/internal/scanner"
)

var show bool

var (
	alertBg = color.New(color.BgRed, color.FgBlack).SprintFunc()
	clearBg = color.New(color.BgGreen, color.FgBlack).SprintFunc()
)

var rootCmd = &cobra.Command{
	Use:   "secr-cli",
	Short: "Scan git repo for secrets in staged, unstaged and working changes",
	Run: func(cobraCmd *cobra.Command, args []string) {
		ruleSet := rules.LoadRules()
		findings, err := scanner.ScanRepo(ruleSet)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error during scan: %v\n", err)
			os.Exit(1)
		}

		if len(findings) > 0 {
			fmt.Println(alertBg(" Potential secrets detected! "))
			
			grouped := make(map[string][]scanner.Finding)
			for _, f := range findings {
				grouped[f.Type] = append(grouped[f.Type], f)
			}

			for _, t := range []string{"staged", "unstaged", "working"} {
				if finds, ok := grouped[t]; ok {
					fmt.Printf("\n%s changes:\n", strings.Title(t))
					for _, f := range finds {
    					fmt.Printf(
    					    "  File: %s\n  Line: %d\n  Rule: %s\n",
    					    f.File, f.Line, f.RuleName,
    					)
    					if show {
    					    fmt.Printf("  Content: %s\n\n", f.Content)
    					} else {
							fmt.Println()
						}
					}
				}
			}
			os.Exit(2)
		} else {
			fmt.Println(clearBg(" No secrets detected! "))
		}
	},
}

func Execute() {
	rootCmd.Flags().BoolVarP(&show, "show", "s", false, "Display secret content in output")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}