package output

import (
	"fmt"
	"strings"

	"github.com/phillips-jon/preflight/internal/checks"
)

// Colors
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

type HumanOutputter struct{}

func (h HumanOutputter) Output(projectName string, results []checks.CheckResult) {
	// Header
	fmt.Println()
	fmt.Printf("%s%s ✈  Preflight Scan Results%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%s   Project: %s%s\n", colorGray, projectName, colorReset)
	fmt.Println()

	// Category icons
	categoryIcons := map[string]string{
		"ENV":       "📋",
		"HEALTH":    "💓",
		"PAYMENTS":  "💳",
		"ERRORS":    "🐛",
		"ANALYTICS": "📊",
		"INFRA":     "🔧",
		"JOBS":      "⚡",
		"SEO":       "🔍",
		"SECURITY":  "🔒",
		"SECRETS":   "🔑",
		"AI":        "🤖",
		"EMAIL":     "📧",
		"AUTH":      "🔐",
		"STORAGE":   "📦",
		"SEARCH":    "🔎",
		"COMM":      "💬",
	}

	// Map check IDs to display categories
	categoryMap := map[string]string{
		"envParity":       "ENV",
		"healthEndpoint":  "HEALTH",
		"stripe":          "PAYMENTS",
		"sentry":          "ERRORS",
		"bugsnag":         "ERRORS",
		"rollbar":         "ERRORS",
		"honeybadger":     "ERRORS",
		"datadog":         "ERRORS",
		"newrelic":        "ERRORS",
		"logrocket":       "ERRORS",
		"plausible":       "ANALYTICS",
		"fathom":          "ANALYTICS",
		"googleAnalytics": "ANALYTICS",
		"mixpanel":        "ANALYTICS",
		"amplitude":       "ANALYTICS",
		"segment":         "ANALYTICS",
		"hotjar":          "ANALYTICS",
		"redis":           "INFRA",
		"sidekiq":         "JOBS",
		"rabbitmq":        "JOBS",
		"seoMeta":         "SEO",
		"securityHeaders": "SECURITY",
		"secrets":         "SECRETS",
		"openai":          "AI",
		"anthropic":       "AI",
		"google_ai":       "AI",
		"auth0":           "AUTH",
		"clerk":           "AUTH",
		"firebase":        "AUTH",
		"supabase":        "AUTH",
		"postmark":        "EMAIL",
		"sendgrid":        "EMAIL",
		"mailgun":         "EMAIL",
		"aws_ses":         "EMAIL",
		"resend":          "EMAIL",
		"aws_s3":          "STORAGE",
		"cloudinary":      "STORAGE",
		"algolia":         "SEARCH",
		"elasticsearch":   "SEARCH",
		"slack":           "COMM",
		"discord":         "COMM",
		"twilio":          "COMM",
		"intercom":        "COMM",
		"crisp":           "COMM",
	}

	// Print results
	for _, r := range results {
		category := categoryMap[r.ID]
		if category == "" {
			category = strings.ToUpper(r.ID)
		}

		icon := categoryIcons[category]
		if icon == "" {
			icon = "•"
		}

		status := formatStatus(r)
		categoryLabel := fmt.Sprintf("%s %-10s", icon, category)

		fmt.Printf("  %s %s%-38s%s %s\n", categoryLabel, colorReset, r.Title, colorReset, status)

		if !r.Passed && r.Message != "" {
			fmt.Printf("  %s                  └─ %s%s\n", colorGray, r.Message, colorReset)
		}
	}

	// Summary
	summary := CalculateSummary(results)
	fmt.Println()
	fmt.Printf("  %s────────────────────────────────────────────────────────%s\n", colorGray, colorReset)
	fmt.Println()

	// Summary with icons
	fmt.Printf("  %s✓ Passed:%s  %s%d%s", colorGreen, colorReset, colorBold, summary.OK, colorReset)
	if summary.Warn > 0 {
		fmt.Printf("    %s⚠ Warnings:%s %s%d%s", colorYellow, colorReset, colorBold, summary.Warn, colorReset)
	}
	if summary.Fail > 0 {
		fmt.Printf("    %s✗ Failed:%s  %s%d%s", colorRed, colorReset, colorBold, summary.Fail, colorReset)
	}
	fmt.Println()
	fmt.Println()

	// Final verdict
	if summary.Fail > 0 {
		fmt.Printf("  %s%s✗ Not ready for launch%s\n", colorBold, colorRed, colorReset)
	} else if summary.Warn > 0 {
		fmt.Printf("  %s%s⚠ Review warnings before launch%s\n", colorBold, colorYellow, colorReset)
	} else {
		fmt.Printf("  %s%s✓ Ready for launch!%s\n", colorBold, colorGreen, colorReset)
	}
	fmt.Println()
}

func formatStatus(r checks.CheckResult) string {
	if r.Passed {
		return fmt.Sprintf("%s%s✓ OK%s", colorBold, colorGreen, colorReset)
	}

	switch r.Severity {
	case checks.SeverityError:
		return fmt.Sprintf("%s%s✗ FAIL%s", colorBold, colorRed, colorReset)
	case checks.SeverityWarn:
		return fmt.Sprintf("%s%s⚠ WARN%s", colorBold, colorYellow, colorReset)
	default:
		return fmt.Sprintf("%s%s⚠ WARN%s", colorBold, colorYellow, colorReset)
	}
}
