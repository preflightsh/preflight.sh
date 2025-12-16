package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var ignoreCmd = &cobra.Command{
	Use:   "ignore <check-id>",
	Short: "Add a check to the ignore list",
	Long: `Add a check ID to the ignore list in preflight.yml.
The check will be skipped in future scans.

Example:
  preflight ignore sitemap
  preflight ignore llmsTxt
  preflight ignore debug_statements`,
	Args: cobra.ExactArgs(1),
	RunE: runIgnore,
}

func init() {
	rootCmd.AddCommand(ignoreCmd)
}

func runIgnore(cmd *cobra.Command, args []string) error {
	checkID := args[0]

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	configPath := filepath.Join(cwd, "preflight.yml")

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("preflight.yml not found. Run 'preflight init' first")
		}
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Parse as generic map to preserve structure
	var cfg map[string]interface{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse preflight.yml: %w", err)
	}

	// Get or create ignore list
	var ignoreList []string
	if existing, ok := cfg["ignore"]; ok {
		if list, ok := existing.([]interface{}); ok {
			for _, item := range list {
				if s, ok := item.(string); ok {
					ignoreList = append(ignoreList, s)
				}
			}
		}
	}

	// Check if already ignored
	for _, id := range ignoreList {
		if id == checkID {
			fmt.Printf("'%s' is already in the ignore list\n", checkID)
			return nil
		}
	}

	// Add to ignore list
	ignoreList = append(ignoreList, checkID)
	cfg["ignore"] = ignoreList

	// Write back
	newData, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	if err := os.WriteFile(configPath, newData, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("Added '%s' to ignore list\n", checkID)
	return nil
}

// Also add an unignore command
var unignoreCmd = &cobra.Command{
	Use:   "unignore <check-id>",
	Short: "Remove a check from the ignore list",
	Long: `Remove a check ID from the ignore list in preflight.yml.

Example:
  preflight unignore sitemap`,
	Args: cobra.ExactArgs(1),
	RunE: runUnignore,
}

func init() {
	rootCmd.AddCommand(unignoreCmd)
}

func runUnignore(cmd *cobra.Command, args []string) error {
	checkID := args[0]

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	configPath := filepath.Join(cwd, "preflight.yml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("preflight.yml not found. Run 'preflight init' first")
		}
		return fmt.Errorf("failed to read config: %w", err)
	}

	var cfg map[string]interface{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse preflight.yml: %w", err)
	}

	// Get ignore list
	var ignoreList []string
	if existing, ok := cfg["ignore"]; ok {
		if list, ok := existing.([]interface{}); ok {
			for _, item := range list {
				if s, ok := item.(string); ok {
					ignoreList = append(ignoreList, s)
				}
			}
		}
	}

	// Find and remove
	found := false
	var newList []string
	for _, id := range ignoreList {
		if id == checkID {
			found = true
		} else {
			newList = append(newList, id)
		}
	}

	if !found {
		fmt.Printf("'%s' is not in the ignore list\n", checkID)
		return nil
	}

	// Update or remove ignore key
	if len(newList) > 0 {
		cfg["ignore"] = newList
	} else {
		delete(cfg, "ignore")
	}

	newData, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	if err := os.WriteFile(configPath, newData, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("Removed '%s' from ignore list\n", checkID)
	return nil
}

// Helper to list available check IDs
var listChecksCmd = &cobra.Command{
	Use:   "checks",
	Short: "List all available check IDs",
	RunE: func(cmd *cobra.Command, args []string) error {
		checkIDs := []string{
			"seoMeta", "canonical", "structured_data", "indexNow (opt-in)",
			"ogTwitter", "viewport", "lang",
			"securityHeaders", "ssl", "www_redirect", "email_auth (opt-in)", "secrets",
			"envParity", "healthEndpoint",
			"sentry", "plausible", "fathom", "googleAnalytics", "redis", "sidekiq", "stripe",
			"vulnerability", "debug_statements", "error_pages", "image_optimization",
			"legal_pages", "cookie_consent (opt-in)",
			"favicon", "robotsTxt", "sitemap", "llmsTxt", "adsTxt (opt-in)", "humansTxt (opt-in)", "license (opt-in)",
		}

		fmt.Println("Available check IDs:")
		for _, id := range checkIDs {
			fmt.Printf("  - %s\n", id)
		}
		fmt.Println()
		fmt.Println("Use 'preflight ignore <check-id>' to silence a check")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listChecksCmd)
}
