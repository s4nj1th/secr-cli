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
		{"GitHub Token", regexp.MustCompile(`(ghp|gho|ghu|ghs|ghr)_[a-zA-Z0-9]{36}`)},
        {"Google API Key", regexp.MustCompile(`AIza[0-9A-Za-z\-_]{35}`)},
        {"Heroku API Key", regexp.MustCompile(`[hH][eE][rR][oO][kK][uU].{0,30}[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}`)},
        {"SSH Private Key", regexp.MustCompile(`-----BEGIN OPENSSH PRIVATE KEY-----`)},
        {"PGP Private Key", regexp.MustCompile(`-----BEGIN PGP PRIVATE KEY BLOCK-----`)},
        {"JWT Token", regexp.MustCompile(`eyJ[a-zA-Z0-9\/_-]{10,}\.[a-zA-Z0-9\/_-]{10,}\.[a-zA-Z0-9\/_-]{10,}`)},
        {"Credit Card Number", regexp.MustCompile(`\b(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|6(?:011|5[0-9]{2})[0-9]{12}|(?:2131|1800|35\d{3})\d{11})\b`)},
        {"Basic Auth Credentials", regexp.MustCompile(`(?i)basic [a-z0-9=:_\+\/-]{5,100}`)},
        {"Docker Registry Auth", regexp.MustCompile(`"auth"\s*:\s*"[a-z0-9=:_\+\/-]{5,100}"`)},
        {"Azure Storage Key", regexp.MustCompile(`DefaultEndpointsProtocol=https;AccountName=[a-z0-9]{3,24};AccountKey=[a-z0-9\/+]{88}==`)},
        {"Google OAuth Token", regexp.MustCompile(`ya29\.[a-zA-Z0-9\-_]+`)},
        {"Facebook Access Token", regexp.MustCompile(`EAACEdEose0cBA[0-9A-Za-z]+`)},
        {"Twitter API Key", regexp.MustCompile(`(?i)twitter(.{0,20})?['"][0-9a-z]{35,44}['"]`)},
        {"Database Connection String", regexp.MustCompile(`(?i)(jdbc:|mongodb:\/\/|postgresql:\/\/|mysql:\/\/).+:[^@]+@[a-z0-9\.-]+`)},
	}
}
