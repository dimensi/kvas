package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var apiURL = "https://api.github.com"

// Release represents minimal information about a GitHub release.
type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// Latest fetches the latest release metadata for the given repository.
// owner and repo correspond to the GitHub repository path.
func Latest(owner, repo string) (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", apiURL, owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var r Release
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

// Install downloads the first asset from the release into dir.
// The downloaded file path is returned. It is marked executable.
func Install(rel *Release, dir string) (string, error) {
	if len(rel.Assets) == 0 {
		return "", fmt.Errorf("no assets in release")
	}
	asset := rel.Assets[0]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	dest := filepath.Join(dir, asset.Name)
	resp, err := http.Get(asset.BrowserDownloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: %s", resp.Status)
	}
	f, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", err
	}
	if err := os.Chmod(dest, 0755); err != nil {
		return "", err
	}
	return dest, nil
}

// Update downloads and installs the latest release asset into dir.
func Update(owner, repo, dir string) (string, error) {
	rel, err := Latest(owner, repo)
	if err != nil {
		return "", err
	}
	return Install(rel, dir)
}
