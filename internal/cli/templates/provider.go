package templates

import (
	"errors"
)

// TemplateProvider defines the interface for obtaining project templates
type TemplateProvider interface {
	// GetTemplate returns the template content for the given version
	GetTemplate(version string) ([]byte, error)

	// DownloadTemplate downloads and extracts a template release to a temporary directory
	DownloadTemplate(releaseTag string) (string, error)

	// ExtractTemplate extracts a template archive to the destination directory
	ExtractTemplate(archivePath, destDir string) error
}

// ErrTemplateNotFound is returned when a template cannot be found
var ErrTemplateNotFound = errors.New("template not found")

// ErrInvalidTemplate is returned when template content is invalid
var ErrInvalidTemplate = errors.New("invalid template content")