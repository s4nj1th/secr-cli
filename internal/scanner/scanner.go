package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"os"
	"io/ioutil"
	"os/exec"
	"strings"

	"secr-cli/internal/rules"
)

type Finding struct {
	File     string
	Line     int
	Content  string
	RuleName string
}

func ScanRepo(staged bool, ruleSet []rules.Rule) ([]Finding, error) {
	if err := checkGitRepo(); err != nil {
		return nil, err
	}

	if staged {
		output, err := gitDiff("--cached")
		if err != nil {
			return nil, fmt.Errorf("git diff failed: %w", err)
		}
		return scanGitDiff(output, ruleSet), nil
	}

	headCmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
	if err := headCmd.Run(); err != nil {
		return scanWorkingDirectory(ruleSet)
	}

	output, err := gitShow()
	if err != nil {
		return nil, fmt.Errorf("git show failed: %w", err)
	}
	return scanGitDiff(output, ruleSet), nil
}


func checkGitRepo() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = wd
	out, err := cmd.Output()
	if err != nil || strings.TrimSpace(string(out)) != "true" {
		return fmt.Errorf("not a git repository")
	}
	return nil
}

func gitShow() ([]byte, error) {
	cmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("no commits to scan (HEAD does not exist)")
	}

	cmd = exec.Command("git", "show", "HEAD")
	cmd.Dir, _ = os.Getwd()
	return cmd.Output()
}


func gitDiff(args ...string) ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("git", append([]string{"diff"}, args...)...)
	cmd.Dir = wd
	return cmd.Output()
}

func scanGitDiff(diff []byte, rules []rules.Rule) []Finding {
	var findings []Finding

	scanner := bufio.NewScanner(bytes.NewReader(diff))
	var currentFile string
	var lineNum int

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "+++ b/") {
			currentFile = strings.TrimPrefix(line, "+++ b/")
			lineNum = 0
			continue
		}
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			lineNum++
			content := line[1:]

			for _, rule := range rules {
				if rule.Pattern.MatchString(content) {
					findings = append(findings, Finding{
						File:     currentFile,
						Line:     lineNum,
						Content:  content,
						RuleName: rule.Name,
					})
				}
			}
		}
	}
	return findings
}

func scanWorkingDirectory(ruleSet []rules.Rule) ([]Finding, error) {
	var findings []Finding

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}

		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			for _, rule := range ruleSet {
				if rule.Pattern.MatchString(line) {
					findings = append(findings, Finding{
						File:     path,
						Line:     i + 1,
						Content:  line,
						RuleName: rule.Name,
					})
				}
			}
		}
		return nil
	})

	return findings, err
}

