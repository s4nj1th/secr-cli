package rules

import "regexp"

type Rule struct {
	Name    string
	Pattern *regexp.Regexp
}

func LoadRules() []Rule {
	return []Rule{
		{"AWS Access Key", regexp.MustCompile(`AKIA[0-9A-Z]{16}`)},
		{"AWS Secret Key", regexp.MustCompile(`(?i)aws(.{0,20})?(secret)?(.{0,20})?['"][0-9a-zA-Z\/+]{40}['"]`)},
		{"Private Key", regexp.MustCompile(`-----BEGIN (RSA|EC|DSA)? PRIVATE KEY-----`)},
		{"Generic API Key", regexp.MustCompile(`(?i)(api|token|secret)[\s:=]+['"]?[a-z0-9]{32,}['"]?`)},
		{"Slack Token", regexp.MustCompile(`xox[baprs]-[0-9a-zA-Z]{10,48}`)},
	}
}
