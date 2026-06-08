package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed project_template
var projectTemplateFS embed.FS

const gitignoreContent = `# Forge generated project files
.env
*.log
bin/
dist/
tmp/
`

// CopyEmbeddedTemplate copies the bundled community template into destDir.
func CopyEmbeddedTemplate(destDir string) error {
	return fs.WalkDir(projectTemplateFS, "project_template", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "project_template" {
			return nil
		}

		relativePath, err := filepath.Rel("project_template", path)
		if err != nil {
			return fmt.Errorf("compute relative path for %s: %w", path, err)
		}
		if relativePath == "go.mod.tmpl" {
			relativePath = "go.mod"
		}

		targetPath := filepath.Join(destDir, relativePath)
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		content, err := projectTemplateFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read embedded template file %s: %w", path, err)
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("create target directory for %s: %w", targetPath, err)
		}

		return os.WriteFile(targetPath, content, 0644)
	})
}

// WriteProjectGitignore writes the standard generated-project .gitignore file.
func WriteProjectGitignore(destDir string) error {
	return os.WriteFile(filepath.Join(destDir, ".gitignore"), []byte(gitignoreContent), 0644)
}
