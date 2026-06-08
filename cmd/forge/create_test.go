package forge

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopyEmbeddedTemplate(t *testing.T) {
	tempDir := t.TempDir()

	if err := CopyEmbeddedTemplate(tempDir); err != nil {
		t.Fatalf("copy embedded template: %v", err)
	}
	if err := WriteProjectGitignore(tempDir); err != nil {
		t.Fatalf("write .gitignore: %v", err)
	}
	if err := GenerateGoMod(tempDir, "github.com/test/project"); err != nil {
		t.Fatalf("generate go.mod: %v", err)
	}

	requiredFiles := []string{
		"go.mod",
		".gitignore",
		filepath.Join("cmd", "api", "main.go"),
		filepath.Join("cmd", "worker", "main.go"),
		filepath.Join("internal", "app", "bootstrap.go"),
		filepath.Join("internal", "config", "config.go"),
	}

	for _, path := range requiredFiles {
		if _, err := os.Stat(filepath.Join(tempDir, path)); err != nil {
			t.Fatalf("expected template file %s: %v", path, err)
		}
	}
}

func TestEmbeddedTemplateUsesModulePlaceholder(t *testing.T) {
	tempDir := t.TempDir()

	if err := CopyEmbeddedTemplate(tempDir); err != nil {
		t.Fatalf("copy embedded template: %v", err)
	}
	if err := WriteProjectGitignore(tempDir); err != nil {
		t.Fatalf("write .gitignore: %v", err)
	}
	if err := GenerateGoMod(tempDir, "{{MODULE_PATH}}"); err != nil {
		t.Fatalf("generate go.mod: %v", err)
	}
	goModBytes, err := os.ReadFile(filepath.Join(tempDir, "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}

	goMod := string(goModBytes)
	if !strings.Contains(goMod, "module") {
		t.Fatalf("expected generated go.mod to contain module declaration, got: %s", goMod)
	}
}

func TestRunCreateGeneratesProjectFromEmbeddedTemplate(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir temp dir: %v", err)
	}
	defer func() {
		_ = os.Chdir(originalWD)
	}()

	projectName := "my_service"
	if err := runCreate(t.Context(), projectName, "", "", "minio", true); err != nil {
		t.Fatalf("runCreate: %v", err)
	}

	goModBytes, err := os.ReadFile(filepath.Join(tempDir, projectName, "go.mod"))
	if err != nil {
		t.Fatalf("read generated go.mod: %v", err)
	}

	goMod := string(goModBytes)
	if strings.Contains(goMod, "{{MODULE_PATH}}") {
		t.Fatalf("expected generated go.mod placeholders to be replaced, got: %s", goMod)
	}
	if !strings.Contains(goMod, "module my_service") {
		t.Fatalf("expected generated go.mod to use project module path, got: %s", goMod)
	}

	gitignoreBytes, err := os.ReadFile(filepath.Join(tempDir, projectName, ".gitignore"))
	if err != nil {
		t.Fatalf("read generated .gitignore: %v", err)
	}

	if !strings.Contains(string(gitignoreBytes), ".env") {
		t.Fatalf("expected generated .gitignore to include env ignore rules, got: %s", string(gitignoreBytes))
	}
}
