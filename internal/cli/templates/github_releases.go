package templates

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// GitHubReleasesProvider implements TemplateProvider using GitHub Releases
type GitHubReleasesProvider struct {
	Owner      string
	Repo       string
	HTTPClient *http.Client
}

// NewGitHubReleasesProvider creates a new GitHub releases template provider
func NewGitHubReleasesProvider(owner, repo string) *GitHubReleasesProvider {
	return &GitHubReleasesProvider{
		Owner:      owner,
		Repo:       repo,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// GetTemplate returns the template content for the given version
// Note: This implementation downloads the release each time. In practice,
// you might want to cache this locally.
func (p *GitHubReleasesProvider) GetTemplate(version string) ([]byte, error) {
    releaseTag := fmt.Sprintf("v%s", version)
    extractedPath, err := p.DownloadTemplate(releaseTag)
    if err != nil {
        return nil, err
    }
    // Expect a file named "template.yaml" (or any file) inside the extracted directory.
    // Try to read the first regular file in the directory.
    entries, err := os.ReadDir(extractedPath)
    if err != nil {
        return nil, fmt.Errorf("reading extracted directory: %w", err)
    }
    for _, entry := range entries {
        if !entry.IsDir() {
            data, err := os.ReadFile(filepath.Join(extractedPath, entry.Name()))
            if err != nil {
                return nil, fmt.Errorf("reading template file %s: %w", entry.Name(), err)
            }
            return data, nil
        }
    }
    return nil, fmt.Errorf("no template file found in %s", extractedPath)
}

// DownloadTemplate downloads and extracts a template release to a temporary directory
func (p *GitHubReleasesProvider) DownloadTemplate(releaseTag string) (string, error) {
	// Get release information from GitHub API
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", p.Owner, p.Repo, releaseTag)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	// Set user agent to avoid GitHub API rate limiting for unauthenticated requests
	req.Header.Set("User-Agent", "Forge-CLI")

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("get release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("%w: release %s not found", ErrTemplateNotFound, releaseTag)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d when fetching release", resp.StatusCode)
	}

	// Parse release data
	var releaseData struct {
		Assets []struct {
			BrowserDownloadURL string `json:"browser_download_url"`
			Name               string `json:"name"`
		} `json:"assets"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&releaseData); err != nil {
		return "", fmt.Errorf("decode release data: %w", err)
	}

	// Find the tarball asset (prefer .tar.gz)
	var downloadURL string
	for _, asset := range releaseData.Assets {
		if asset.Name == "source.tar.gz" || filepath.Ext(asset.Name) == ".tar.gz" {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return "", fmt.Errorf("no tarball asset found in release %s", releaseTag)
	}

	// Download the tarball
	tmpDir, err := os.MkdirTemp("", "forge-github-*")
	if err != nil {
		return "", fmt.Errorf("create temp dir: %w", err)
	}

	tarballPath := filepath.Join(tmpDir, "release.tar.gz")
	resp, err = p.HTTPClient.Get(downloadURL)
	if err != nil {
		return "", fmt.Errorf("download tarball: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download tarball: %s", resp.Status)
	}

	out, err := os.Create(tarballPath)
	if err != nil {
		return "", fmt.Errorf("create tarball file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("save tarball: %w", err)
	}

	// Extract the tarball
	extractedPath := filepath.Join(tmpDir, "extracted")
	if err := os.MkdirAll(extractedPath, 0755); err != nil {
		return "", fmt.Errorf("create extract dir: %w", err)
	}

	cmd := exec.Command("tar", "-xzf", tarballPath, "-C", extractedPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("extract tarball: %w, output: %s", err, string(output))
	}

	// Find the extracted directory (GitHub releases typically create a dir like owner-repo-tag)
	entries, err := os.ReadDir(extractedPath)
	if err != nil {
		return "", fmt.Errorf("read extracted dir: %w", err)
	}
	if len(entries) == 0 {
		return "", fmt.Errorf("no content extracted from release")
	}

	// Return the path to the extracted content
	return filepath.Join(extractedPath, entries[0].Name()), nil
}

// ExtractTemplate extracts a template archive to the destination directory
// Note: For GitHub releases, we expect the archive to already be extracted
// by DownloadTemplate, so this is a no-op that just validates the directory exists.
func (p *GitHubReleasesProvider) ExtractTemplate(archivePath, destDir string) error {
	// In this implementation, DownloadTemplate already extracts the template
	// so we just verify the source directory exists and copy it
	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		return fmt.Errorf("archive path does not exist: %s", archivePath)
	}

	// Copy the extracted template to destination
	return copyDirectory(archivePath, destDir)
}

// copyDirectory recursively copies a directory tree
func copyDirectory(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("read source directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, entry.Type().Perm()); err != nil {
				return fmt.Errorf("create directory %s: %w", dstPath, err)
			}
			if err := copyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return fmt.Errorf("copy file %s: %w", srcPath, err)
			}
		}
	}

	return nil
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create destination file: %w", err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}