package forge

import "testing"

func TestBuildTemplateCandidateURLsCommunity(t *testing.T) {
	urls := buildTemplateCandidateURLs("1.0.0")

	expected := []string{
		"https://github.com/AARCSX/AARCSX_Forge/releases/download/1.0.0/forge-template-1.0.0.tar.gz",
		"https://github.com/AARCSX/AARCSX_Forge/releases/download/v1.0.0/forge-template-1.0.0.tar.gz",
		"https://github.com/AARCSX/AARCSX_Forge/releases/download/1.0.1/forge-template-1.0.1.tar.gz",
		"https://github.com/AARCSX/AARCSX_Forge/releases/download/v1.0.1/forge-template-1.0.1.tar.gz",
	}

	if len(urls) != len(expected) {
		t.Fatalf("expected %d URLs, got %d: %#v", len(expected), len(urls), urls)
	}

	for index, want := range expected {
		if urls[index] != want {
			t.Fatalf("expected URL %d to be %q, got %q", index, want, urls[index])
		}
	}
}

func TestBuildTemplateCandidateURLsDeduplicatesVersions(t *testing.T) {
	urls := buildTemplateCandidateURLs("1.0.1")

	expected := []string{
		"https://github.com/AARCSX/AARCSX_Forge/releases/download/1.0.1/forge-template-1.0.1.tar.gz",
		"https://github.com/AARCSX/AARCSX_Forge/releases/download/v1.0.1/forge-template-1.0.1.tar.gz",
	}

	if len(urls) != len(expected) {
		t.Fatalf("expected %d URLs, got %d: %#v", len(expected), len(urls), urls)
	}

	for index, want := range expected {
		if urls[index] != want {
			t.Fatalf("expected URL %d to be %q, got %q", index, want, urls[index])
		}
	}
}
