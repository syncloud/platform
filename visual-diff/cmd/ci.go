package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type build struct {
	Number int    `json:"number"`
	Source string `json:"source"`
	Status string `json:"status"`
}

type fileEntry struct {
	Name string `json:"name"`
}

type CIClient struct {
	apiURL   string
	filesURL string
	arch     string
	http     *http.Client
}

func NewCIClient(apiURL, filesURL, arch string, httpClient *http.Client) *CIClient {
	return &CIClient{apiURL: apiURL, filesURL: filesURL, arch: arch, http: httpClient}
}

func (c *CIClient) FindLatestBuild(branch string) string {
	url := fmt.Sprintf("%s/builds?limit=10&branch=%s", c.apiURL, branch)
	resp, err := c.http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var builds []build
	if err := json.NewDecoder(resp.Body).Decode(&builds); err != nil {
		return ""
	}
	for _, b := range builds {
		if b.Status == "success" {
			return fmt.Sprintf("%d", b.Number)
		}
	}
	return ""
}

func (c *CIClient) DownloadBranch(branch, outDir string) error {
	fmt.Printf("Finding latest build for '%s'...\n", branch)
	b := c.FindLatestBuild(branch)
	if b == "" {
		return fmt.Errorf("no successful build for branch %q", branch)
	}
	fmt.Printf("Branch: %s (build #%s)\n", branch, b)
	return c.DownloadBuild(b, outDir)
}

func (c *CIClient) DownloadBuild(buildNum, outDir string) error {
	os.RemoveAll(outDir)
	for _, view := range []string{"desktop", "mobile"} {
		url := fmt.Sprintf("%s/%s-%s/distro/%s/screenshot", c.filesURL, buildNum, c.arch, view)
		dir := filepath.Join(outDir, view)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		fmt.Printf("\n=== %s screenshots ===\n", view)
		if err := c.downloadPNGs(url, dir); err != nil {
			return err
		}
	}
	return nil
}

func (c *CIClient) downloadPNGs(url, dir string) error {
	names, err := c.listRemotePNGs(url)
	if err != nil {
		return err
	}
	for _, name := range names {
		target := filepath.Join(dir, name)
		if _, err := os.Stat(target); err == nil {
			continue
		}
		if err := c.downloadFile(url+"/"+name, target); err != nil {
			return fmt.Errorf("downloading %s: %w", name, err)
		}
		fmt.Printf("  %s\n", name)
	}
	return nil
}

func (c *CIClient) listRemotePNGs(url string) ([]string, error) {
	resp, err := c.http.Get(url + "/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var files []fileEntry
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, nil
	}
	var names []string
	for _, f := range files {
		if strings.HasSuffix(f.Name, ".png") {
			names = append(names, f.Name)
		}
	}
	return names, nil
}

func (c *CIClient) downloadFile(url, path string) error {
	resp, err := c.http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}
