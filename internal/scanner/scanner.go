package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"secr-cli/internal/gitignore"
	"secr-cli/internal/rules"
)

type Finding struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Content  string `json:"content"`
	RuleName string `json:"rule"`
	Severity string `json:"severity"`
	Type     string `json:"type"`
}

type ScanOptions struct {
	StagedOnly  bool
	NoGitignore bool
	Workers     int
}

func ScanRepo(ruleSet []rules.Rule, opts ScanOptions) ([]Finding, error) {
	if err := checkGitRepo(); err != nil {
		return nil, err
	}

	workers := opts.Workers
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	var allFindings []Finding
	var mu sync.Mutex

	var wg sync.WaitGroup

	wg.Add(1)
	var stagedErr error
	go func() {
		defer wg.Done()
		output, err := gitDiff("--cached")
		if err != nil {
			stagedErr = fmt.Errorf("git diff (staged) failed: %w", err)
			return
		}
		findings := scanGitDiff(output, ruleSet, "staged")
		mu.Lock()
		allFindings = append(allFindings, findings...)
		mu.Unlock()
	}()

	if !opts.StagedOnly {
		wg.Add(1)
		var unstagedErr error
		go func() {
			defer wg.Done()
			output, err := gitDiff()
			if err != nil {
				unstagedErr = fmt.Errorf("git diff (unstaged) failed: %w", err)
				return
			}
			findings := scanGitDiff(output, ruleSet, "unstaged")
			mu.Lock()
			allFindings = append(allFindings, findings...)
			mu.Unlock()
			_ = unstagedErr
		}()
	}

	wg.Wait()

	if stagedErr != nil {
		return nil, stagedErr
	}

	if !opts.StagedOnly {
		workingFindings, err := scanWorkingDirectory(ruleSet, opts.NoGitignore, workers)
		if err != nil {
			return nil, fmt.Errorf("working directory scan failed: %w", err)
		}
		allFindings = append(allFindings, workingFindings...)
	}

	return allFindings, nil
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

func gitDiff(args ...string) ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("git", append([]string{"diff"}, args...)...)
	cmd.Dir = wd
	return cmd.Output()
}

func scanGitDiff(diff []byte, ruleSet []rules.Rule, diffType string) []Finding {
	var findings []Finding

	sc := bufio.NewScanner(bytes.NewReader(diff))
	var currentFile string
	var lineNum int

	for sc.Scan() {
		line := sc.Text()

		if strings.HasPrefix(line, "+++ b/") {
			currentFile = strings.TrimPrefix(line, "+++ b/")
			lineNum = 0
			continue
		}
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			lineNum++
			content := line[1:]

			for _, rule := range ruleSet {
				if rule.Pattern.MatchString(content) {
					findings = append(findings, Finding{
						File:     currentFile,
						Line:     lineNum,
						Content:  content,
						RuleName: rule.Name,
						Severity: string(rule.Severity),
						Type:     diffType,
					})
				}
			}
		}
	}
	return findings
}

func scanWorkingDirectory(ruleSet []rules.Rule, noGitignore bool, workers int) ([]Finding, error) {
	var files []string
	var err error

	if noGitignore {
		files, err = gitignore.AllFiles()
	} else {
		files, err = gitignore.TrackedFiles()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	fileCh := make(chan string, len(files))
	findingCh := make(chan []Finding, len(files))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range fileCh {
				results := scanFile(path, ruleSet)
				if len(results) > 0 {
					findingCh <- results
				}
			}
		}()
	}

	for _, f := range files {
		fileCh <- f
	}
	close(fileCh)

	go func() {
		wg.Wait()
		close(findingCh)
	}()

	var allFindings []Finding
	for batch := range findingCh {
		allFindings = append(allFindings, batch...)
	}

	return allFindings, nil
}

func scanFile(path string, ruleSet []rules.Rule) []Finding {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil || isBinary(buf[:n]) {
		return nil
	}

	if _, err := f.Seek(0, 0); err != nil {
		return nil
	}

	var findings []Finding
	sc := bufio.NewScanner(f)
	lineNum := 0

	for sc.Scan() {
		lineNum++
		line := sc.Text()

		for _, rule := range ruleSet {
			if rule.Pattern.MatchString(line) {
				findings = append(findings, Finding{
					File:     path,
					Line:     lineNum,
					Content:  line,
					RuleName: rule.Name,
					Severity: string(rule.Severity),
					Type:     "working",
				})
			}
		}
	}

	return findings
}

func isBinary(data []byte) bool {
	for _, b := range data {
		if b == 0 {
			return true
		}
	}
	return false
}
