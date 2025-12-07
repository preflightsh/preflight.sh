package checks

import (
	"os"
	"path/filepath"
	"regexp"
)

type FaviconCheck struct{}

func (c FaviconCheck) ID() string {
	return "favicon"
}

func (c FaviconCheck) Title() string {
	return "Favicon and app icons present"
}

func (c FaviconCheck) Run(ctx Context) (CheckResult, error) {
	var found []string
	var missing []string

	// Common web root directories across frameworks
	webRoots := []string{
		"public",  // Laravel, Rails, many Node.js
		"static",  // Hugo, some SSGs
		"web",     // Craft CMS, Symfony
		"www",     // Some PHP apps
		"dist",    // Built static sites
		"build",   // Build outputs
		"_site",   // Jekyll
		"out",     // Next.js static export
		"app",     // Next.js App Router
		"",        // Root directory
	}

	// Check for common favicon locations
	faviconFiles := []string{"favicon.ico", "favicon.png", "favicon.svg", "favicon.webp", "icon.png", "icon.svg"}
	var faviconPaths []string
	for _, root := range webRoots {
		for _, file := range faviconFiles {
			if root == "" {
				faviconPaths = append(faviconPaths, file)
			} else {
				faviconPaths = append(faviconPaths, root+"/"+file)
				// Also check assets subdirectories
				faviconPaths = append(faviconPaths, root+"/assets/"+file)
				faviconPaths = append(faviconPaths, root+"/assets/images/"+file)
				faviconPaths = append(faviconPaths, root+"/images/"+file)
				faviconPaths = append(faviconPaths, root+"/img/"+file)
			}
		}
	}

	hasFavicon := false
	for _, path := range faviconPaths {
		fullPath := filepath.Join(ctx.RootDir, path)
		if _, err := os.Stat(fullPath); err == nil {
			hasFavicon = true
			found = append(found, path)
			break
		}
	}

	if !hasFavicon {
		missing = append(missing, "favicon")
	}

	// Check for Apple Touch Icon
	var appleTouchPaths []string
	for _, root := range webRoots {
		if root == "" {
			appleTouchPaths = append(appleTouchPaths, "apple-touch-icon.png", "apple-icon.png")
		} else {
			appleTouchPaths = append(appleTouchPaths,
				root+"/apple-touch-icon.png",
				root+"/apple-icon.png",
			)
		}
	}

	hasAppleIcon := false
	for _, path := range appleTouchPaths {
		fullPath := filepath.Join(ctx.RootDir, path)
		if _, err := os.Stat(fullPath); err == nil {
			hasAppleIcon = true
			found = append(found, path)
			break
		}
	}

	// Also check HTML/templates for apple-touch-icon link
	if !hasAppleIcon {
		// Check configured layout first
		cfg := ctx.Config.Checks.SEOMeta
		if cfg != nil && cfg.MainLayout != "" {
			layoutPath := filepath.Join(ctx.RootDir, cfg.MainLayout)
			if content, err := os.ReadFile(layoutPath); err == nil {
				if regexp.MustCompile(`(?i)apple-touch-icon`).Match(content) {
					hasAppleIcon = true
					found = append(found, "apple-touch-icon (in HTML)")
				}
			}
		}

		// Check common template locations
		if !hasAppleIcon {
			templatePaths := []string{
				"templates/_layout.twig",           // Craft CMS
				"templates/_layout.html",           // Craft CMS
				"templates/_head.twig",             // Craft CMS partials
				"templates/_head.html",
				"templates/_partials/head.twig",    // Craft CMS partials
				"templates/_partials/header.twig",  // Craft CMS partials
				"app/views/layouts/application.html.erb", // Rails
				"resources/views/layouts/app.blade.php",  // Laravel
				"_includes/head.html",              // Jekyll
				"layouts/_default/baseof.html",     // Hugo
				"src/layouts/Layout.astro",         // Astro
			}
			for _, tplPath := range templatePaths {
				fullPath := filepath.Join(ctx.RootDir, tplPath)
				if content, err := os.ReadFile(fullPath); err == nil {
					if regexp.MustCompile(`(?i)apple-touch-icon`).Match(content) {
						hasAppleIcon = true
						found = append(found, "apple-touch-icon (in HTML)")
						break
					}
				}
			}
		}
	}

	if !hasAppleIcon {
		missing = append(missing, "apple-touch-icon")
	}

	// Check for web app manifest
	var manifestPaths []string
	for _, root := range webRoots {
		if root == "" {
			manifestPaths = append(manifestPaths, "manifest.json", "site.webmanifest")
		} else {
			manifestPaths = append(manifestPaths,
				root+"/manifest.json",
				root+"/site.webmanifest",
				root+"/manifest.ts",
				root+"/manifest.js",
			)
		}
	}

	hasManifest := false
	for _, path := range manifestPaths {
		fullPath := filepath.Join(ctx.RootDir, path)
		if _, err := os.Stat(fullPath); err == nil {
			hasManifest = true
			found = append(found, path)
			break
		}
	}

	if !hasManifest {
		missing = append(missing, "web manifest")
	}

	// Determine result
	if len(missing) == 0 {
		return CheckResult{
			ID:       c.ID(),
			Title:    c.Title(),
			Severity: SeverityInfo,
			Passed:   true,
			Message:  "All icons and manifest present",
		}, nil
	}

	if hasFavicon && len(missing) <= 2 {
		// Has favicon but missing apple icon or manifest - just warn
		return CheckResult{
			ID:       c.ID(),
			Title:    c.Title(),
			Severity: SeverityWarn,
			Passed:   false,
			Message:  "Missing: " + joinStrings(missing, ", "),
			Suggestions: []string{
				"Add apple-touch-icon.png (180x180px) for iOS",
				"Add manifest.json for PWA support",
			},
		}, nil
	}

	return CheckResult{
		ID:       c.ID(),
		Title:    c.Title(),
		Severity: SeverityError,
		Passed:   false,
		Message:  "Missing favicon",
		Suggestions: []string{
			"Add favicon.ico or favicon.png to public/",
			"Use https://realfavicongenerator.net for complete icon set",
		},
	}, nil
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
