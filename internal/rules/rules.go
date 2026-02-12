package rules

import "regexp"

type Severity string

const (
	SeverityHigh   Severity = "HIGH"
	SeverityMedium Severity = "MEDIUM"
	SeverityLow    Severity = "LOW"
)

type Rule struct {
	Name     string
	Pattern  *regexp.Regexp
	Severity Severity
}

func LoadRules() []Rule {
	return []Rule{
		// Cloud Credentials
		{"AWS Access Key", regexp.MustCompile(`AKIA[0-9A-Z]{16}`), SeverityHigh},
		{"AWS Secret Key", regexp.MustCompile(`(?i)aws(.{0,20})?(secret)?(.{0,20})?['"][0-9a-zA-Z\/+]{40}['"]`), SeverityHigh},
		{"Google API Key", regexp.MustCompile(`AIza[0-9A-Za-z\-_]{35}`), SeverityHigh},
		{"Google OAuth Token", regexp.MustCompile(`ya29\.[a-zA-Z0-9\-_]+`), SeverityHigh},
		{"Azure Storage Key", regexp.MustCompile(`DefaultEndpointsProtocol=https;AccountName=[a-z0-9]{3,24};AccountKey=[a-z0-9\/+]{88}==`), SeverityHigh},
		{"Heroku API Key", regexp.MustCompile(`[hH][eE][rR][oO][kK][uU].{0,30}[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}`), SeverityHigh},
		{"DigitalOcean Token", regexp.MustCompile(`dop_v1_[a-f0-9]{64}`), SeverityHigh},
		{"Cloudflare API Key", regexp.MustCompile(`(?i)cloudflare(.{0,20})?['"][a-z0-9]{37}['"]`), SeverityHigh},

		// API Tokens
		{"GitHub Token", regexp.MustCompile(`(ghp|gho|ghu|ghs|ghr)_[a-zA-Z0-9]{36}`), SeverityHigh},
		{"Slack Token", regexp.MustCompile(`xox[baprs]-[0-9a-zA-Z]{10,48}`), SeverityHigh},
		{"Stripe Secret Key", regexp.MustCompile(`sk_live_[0-9a-zA-Z]{24,99}`), SeverityHigh},
		{"Stripe Publishable Key", regexp.MustCompile(`pk_live_[0-9a-zA-Z]{24,99}`), SeverityLow},
		{"SendGrid API Key", regexp.MustCompile(`SG\.[a-zA-Z0-9\-_]{22,}\.[a-zA-Z0-9\-_]{22,}`), SeverityHigh},
		{"Twilio API Key", regexp.MustCompile(`SK[a-f0-9]{32}`), SeverityHigh},
		{"npm Access Token", regexp.MustCompile(`npm_[a-zA-Z0-9]{36}`), SeverityHigh},
		{"PyPI API Token", regexp.MustCompile(`pypi-AgEIcHlwaS5vcmc[a-zA-Z0-9\-_]{50,}`), SeverityHigh},
		{"Discord Bot Token", regexp.MustCompile(`[MN][A-Za-z\d]{23,}\.[\w-]{6}\.[\w-]{27}`), SeverityHigh},
		{"Telegram Bot Token", regexp.MustCompile(`[0-9]{8,10}:[a-zA-Z0-9_-]{35}`), SeverityMedium},
		{"Mailgun API Key", regexp.MustCompile(`key-[a-zA-Z0-9]{32}`), SeverityHigh},
		{"Datadog API Key", regexp.MustCompile(`(?i)datadog(.{0,20})?['"][a-f0-9]{32}['"]`), SeverityMedium},
		{"Shopify Token", regexp.MustCompile(`shpat_[a-fA-F0-9]{32}`), SeverityHigh},
		{"Linear API Key", regexp.MustCompile(`lin_api_[a-zA-Z0-9]{40}`), SeverityMedium},
		{"OpenAI API Key", regexp.MustCompile(`sk-[a-zA-Z0-9]{20}T3BlbkFJ[a-zA-Z0-9]{20}`), SeverityHigh},
		{"Anthropic API Key", regexp.MustCompile(`sk-ant-api03-[a-zA-Z0-9\-_]{93}`), SeverityHigh},
		{"Facebook Access Token", regexp.MustCompile(`EAACEdEose0cBA[0-9A-Za-z]+`), SeverityHigh},
		{"Twitter API Key", regexp.MustCompile(`(?i)twitter(.{0,20})?['"][0-9a-z]{35,44}['"]`), SeverityMedium},

		// Cryptographic Material
		{"RSA Private Key", regexp.MustCompile(`-----BEGIN RSA PRIVATE KEY-----`), SeverityHigh},
		{"EC Private Key", regexp.MustCompile(`-----BEGIN EC PRIVATE KEY-----`), SeverityHigh},
		{"DSA Private Key", regexp.MustCompile(`-----BEGIN DSA PRIVATE KEY-----`), SeverityHigh},
		{"SSH Private Key", regexp.MustCompile(`-----BEGIN OPENSSH PRIVATE KEY-----`), SeverityHigh},
		{"PGP Private Key", regexp.MustCompile(`-----BEGIN PGP PRIVATE KEY BLOCK-----`), SeverityHigh},
		{"Generic Private Key", regexp.MustCompile(`-----BEGIN PRIVATE KEY-----`), SeverityHigh},

		// Database & Connection Strings
		{"Database Connection String", regexp.MustCompile(`(?i)(jdbc:|mongodb://|mongodb\+srv://|postgresql://|postgres://|mysql://|redis://).+:[^@]+@[a-z0-9\.\-]+`), SeverityHigh},
		{"Password in URL", regexp.MustCompile(`[a-zA-Z]{3,10}://[^/\s:@]{3,20}:[^/\s:@]{3,20}@.{1,100}`), SeverityHigh},

		// Authentication
		{"JWT Token", regexp.MustCompile(`eyJ[a-zA-Z0-9\/_-]{10,}\.[a-zA-Z0-9\/_-]{10,}\.[a-zA-Z0-9\/_-]{10,}`), SeverityMedium},
		{"Basic Auth Credentials", regexp.MustCompile(`(?i)basic [a-z0-9=:_\+\/-]{5,100}`), SeverityMedium},
		{"Docker Registry Auth", regexp.MustCompile(`"auth"\s*:\s*"[a-z0-9=:_\+\/-]{5,100}"`), SeverityHigh},

		// Generic Patterns
		{"Generic API Key", regexp.MustCompile(`(?i)(api_key|apikey|api-key|secret_key|secret-key|access_key|access_token|auth_token|service_key)[\s:=]+['"]?[a-z0-9]{32,}['"]?`), SeverityMedium},
		{"Env File Secret", regexp.MustCompile(`(?i)(password|passwd|pwd|secret|token|api_key|apikey)\s*=\s*[^\s]{8,}`), SeverityMedium},
		{"Credit Card Number", regexp.MustCompile(`\b(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|6(?:011|5[0-9]{2})[0-9]{12}|(?:2131|1800|35\d{3})\d{11})\b`), SeverityHigh},
	}
}
