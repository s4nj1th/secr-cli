package gitignore

import (
	"os/exec"
	"strings"
)

func TrackedFiles() ([]string, error) {

	tracked, err := gitLsFiles()
	if err != nil {
		return nil, err
	}

	untracked, err := gitLsFilesUntracked()
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{}, len(tracked)+len(untracked))
	var result []string

	for _, f := range tracked {
		if _, ok := seen[f]; !ok {
			seen[f] = struct{}{}
			result = append(result, f)
		}
	}
	for _, f := range untracked {
		if _, ok := seen[f]; !ok {
			seen[f] = struct{}{}
			result = append(result, f)
		}
	}

	return result, nil
}

func AllFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-files", "--cached", "--others")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parseLines(out), nil
}

func gitLsFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-files")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parseLines(out), nil
}

func gitLsFilesUntracked() ([]string, error) {
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parseLines(out), nil
}

func parseLines(data []byte) []string {
	raw := strings.TrimSpace(string(data))
	if raw == "" {
		return nil
	}
	return strings.Split(raw, "\n")
}
