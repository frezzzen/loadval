package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	CurrentVersion = "0.0.2"
	GitHubOwner    = "frezzzen"
	GitHubRepo     = "loadval"
)

type UpdateAPI struct{}

func NewUpdateAPI() *UpdateAPI {
	return &UpdateAPI{}
}

type GitHubRelease struct {
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Body       string `json:"body"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
	Assets     []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Size               int64  `json:"size"`
	} `json:"assets"`
	PublishedAt time.Time `json:"published_at"`
}

type UpdateInfo struct {
	Available      bool   `json:"available"`
	CurrentVersion string `json:"currentVersion"`
	LatestVersion  string `json:"latestVersion"`
	ReleaseNotes   string `json:"releaseNotes"`
	DownloadURL    string `json:"downloadUrl"`
	DownloadSize   int64  `json:"downloadSize"`
	PublishedAt    string `json:"publishedAt"`
}

type DownloadProgress struct {
	BytesDownloaded int64   `json:"bytesDownloaded"`
	TotalBytes      int64   `json:"totalBytes"`
	Percentage      float64 `json:"percentage"`
	Complete        bool    `json:"complete"`
	Error           string  `json:"error"`
}

func (u *UpdateAPI) GetCurrentVersion() string {
	return CurrentVersion
}

func (u *UpdateAPI) CheckForUpdates() (*UpdateInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", GitHubOwner, GitHubRepo)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if release.Draft || release.Prerelease {
		return &UpdateInfo{
			Available:      false,
			CurrentVersion: CurrentVersion,
			LatestVersion:  CurrentVersion,
		}, nil
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	updateAvailable := compareVersions(latestVersion, CurrentVersion) > 0

	updateInfo := &UpdateInfo{
		Available:      updateAvailable,
		CurrentVersion: CurrentVersion,
		LatestVersion:  latestVersion,
		ReleaseNotes:   release.Body,
		PublishedAt:    release.PublishedAt.Format("2006-01-02"),
	}

	if updateAvailable {
		for _, asset := range release.Assets {
			assetName := strings.ToLower(asset.Name)
			if runtime.GOOS == "windows" && (strings.Contains(assetName, "installer.exe") || strings.Contains(assetName, "windows")) {
				updateInfo.DownloadURL = asset.BrowserDownloadURL
				updateInfo.DownloadSize = asset.Size
				break
			} else if runtime.GOOS == "darwin" && (strings.Contains(assetName, ".dmg") || strings.Contains(assetName, "darwin")) {
				updateInfo.DownloadURL = asset.BrowserDownloadURL
				updateInfo.DownloadSize = asset.Size
				break
			} else if runtime.GOOS == "linux" && (strings.Contains(assetName, "linux") || strings.HasSuffix(assetName, ".appimage")) {
				updateInfo.DownloadURL = asset.BrowserDownloadURL
				updateInfo.DownloadSize = asset.Size
				break
			}
		}
	}

	return updateInfo, nil
}

func (u *UpdateAPI) DownloadUpdate(downloadURL string) (string, error) {
	if downloadURL == "" {
		return "", fmt.Errorf("download URL is empty")
	}

	downloadsDir := getDownloadsDir()
	if err := os.MkdirAll(downloadsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create downloads directory: %w", err)
	}

	filename := filepath.Base(downloadURL)
	if idx := strings.Index(filename, "?"); idx != -1 {
		filename = filename[:idx]
	}
	filePath := filepath.Join(downloadsDir, filename)

	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	client := &http.Client{
		Timeout: 0,
	}

	resp, err := client.Get(downloadURL)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}

func (u *UpdateAPI) DownloadUpdateWithProgress(downloadURL string, progressCallback func(DownloadProgress)) (string, error) {
	if downloadURL == "" {
		return "", fmt.Errorf("download URL is empty")
	}

	downloadsDir := getDownloadsDir()
	if err := os.MkdirAll(downloadsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create downloads directory: %w", err)
	}

	filename := filepath.Base(downloadURL)
	if idx := strings.Index(filename, "?"); idx != -1 {
		filename = filename[:idx]
	}
	filePath := filepath.Join(downloadsDir, filename)

	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	client := &http.Client{
		Timeout: 0,
	}

	resp, err := client.Get(downloadURL)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	totalSize := resp.ContentLength
	downloaded := int64(0)

	buffer := make([]byte, 32*1024)
	lastUpdate := time.Now()

	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			if _, writeErr := out.Write(buffer[:n]); writeErr != nil {
				return "", fmt.Errorf("failed to write file: %w", writeErr)
			}
			downloaded += int64(n)

			if time.Since(lastUpdate) >= 100*time.Millisecond {
				percentage := 0.0
				if totalSize > 0 {
					percentage = float64(downloaded) / float64(totalSize) * 100
				}

				if progressCallback != nil {
					progressCallback(DownloadProgress{
						BytesDownloaded: downloaded,
						TotalBytes:      totalSize,
						Percentage:      percentage,
						Complete:        false,
					})
				}
				lastUpdate = time.Now()
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to read response: %w", err)
		}
	}

	if progressCallback != nil {
		progressCallback(DownloadProgress{
			BytesDownloaded: downloaded,
			TotalBytes:      totalSize,
			Percentage:      100.0,
			Complete:        true,
		})
	}

	return filePath, nil
}

func (u *UpdateAPI) InstallUpdate(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("update file does not exist: %s", filePath)
	}

	switch runtime.GOOS {
	case "windows":
		programFiles := os.Getenv("ProgramFiles")
		if programFiles == "" {
			programFiles = "C:\\Program Files"
		}
		exePath := filepath.Join(programFiles, "frezzzen", "LOADVAL", "loadval.exe")

		// Use PowerShell to run the installer and then restart the app
		powerShellScript := fmt.Sprintf(`
Start-Process -FilePath '%s' -ArgumentList '/S' -Verb RunAs -Wait
Start-Sleep -Seconds 2
Start-Process -FilePath '%s'
`, filePath, exePath)

		cmd := exec.Command("powershell", "-Command", powerShellScript)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start installer: %w", err)
		}

		// Exit the current application after a short delay
		go func() {
			time.Sleep(500 * time.Millisecond)
			os.Exit(0)
		}()

		return nil

	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

func (u *UpdateAPI) OpenDownloadFolder() error {
	downloadsDir := getDownloadsDir()

	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("explorer", downloadsDir)
		return cmd.Start()
	case "darwin":
		cmd := exec.Command("open", downloadsDir)
		return cmd.Start()
	case "linux":
		cmd := exec.Command("xdg-open", downloadsDir)
		return cmd.Start()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

func getDownloadsDir() string {
	appData := os.Getenv("LOCALAPPDATA")
	if appData == "" {
		appData = os.Getenv("APPDATA")
	}
	if appData == "" {
		appData = os.TempDir()
	}
	return filepath.Join(appData, "LOADVAL", "Updates")
}

func compareVersions(v1, v2 string) int {
	v1Parts := strings.Split(strings.TrimPrefix(v1, "v"), ".")
	v2Parts := strings.Split(strings.TrimPrefix(v2, "v"), ".")

	maxLen := len(v1Parts)
	if len(v2Parts) > maxLen {
		maxLen = len(v2Parts)
	}

	for i := 0; i < maxLen; i++ {
		var part1, part2 int
		if i < len(v1Parts) {
			fmt.Sscanf(v1Parts[i], "%d", &part1)
		}
		if i < len(v2Parts) {
			fmt.Sscanf(v2Parts[i], "%d", &part2)
		}

		if part1 > part2 {
			return 1
		}
		if part1 < part2 {
			return -1
		}
	}

	return 0
}
