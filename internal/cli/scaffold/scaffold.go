package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReplacePlaceholders replaces all occurrences of placeholders in files under the given directory
func ReplacePlaceholders(dir string, replacements map[string]string) error {
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Skip certain file types or directories if needed
		if shouldSkipFile(path) {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file %s: %w", path, err)
		}

		// Replace placeholders
		newContent := string(content)
		for placeholder, value := range replacements {
			newContent = strings.ReplaceAll(newContent, placeholder, value)
		}

		// Write back if changed
		if newContent != string(content) {
			if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
				return fmt.Errorf("write file %s: %w", path, err)
			}
		}

		return nil
	})
}

// shouldSkipFile determines if a file should be skipped during placeholder replacement
func shouldSkipFile(path string) bool {
	// Skip binary files, images, archives, etc.
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".tiff", ".ico":
		return true
	case ".zip", ".tar", ".gz", ".bz2", ".xz", ".7z":
		return true
	case ".exe", ".dll", ".so", ".dylib":
		return true
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx":
		return true
	}

	// Skip hidden files and directories (except those we specifically want to process)
	if strings.HasPrefix(filepath.Base(path), ".") &&
	   !strings.HasPrefix(filepath.Base(path), ".forge") {
		return true
	}

	return false
}

// ModuleNameFromProjectName converts a project name to a suitable Go module name
func ModuleNameFromProjectName(projectName string) string {
	// Convert to lowercase and replace hyphens with underscores for module name
	return strings.ToLower(strings.ReplaceAll(projectName, "-", "_"))
}