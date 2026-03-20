package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	ciAPI          = "http://ci.syncloud.org:8080/api/repos/syncloud/platform"
	ciFiles        = "http://ci.syncloud.org:8081/files/platform"
	arch           = "amd64"
	pixelThreshold = 100
	fuzzPercent    = 5.0
)

var views = []string{"desktop", "mobile"}

type Build struct {
	Number int    `json:"number"`
	Source string `json:"source"`
	Status string `json:"status"`
}

type FileEntry struct {
	Name string `json:"name"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: visual-diff <command> [args]")
		fmt.Println("Commands:")
		fmt.Println("  diff <base-dir> <compare-dir> [diff-dir]    Compare two local screenshot dirs")
		fmt.Println("  ci-diff <local-dir> [skip-build]            Compare local screenshots against stable")
		fmt.Println("  download <branch> <output-dir>              Download screenshots from CI")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "diff":
		cmdDiff(os.Args[2:])
	case "ci-diff":
		cmdCIDiff(os.Args[2:])
	case "download":
		cmdDownload(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func cmdDiff(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: visual-diff diff <base-dir> <compare-dir> [diff-dir]")
		os.Exit(1)
	}
	baseDir := args[0]
	cmpDir := args[1]
	diffDir := ""
	if len(args) > 2 {
		diffDir = args[2]
	}

	total := 0
	for _, view := range views {
		baseView := filepath.Join(baseDir, view)
		cmpView := filepath.Join(cmpDir, view)
		if !dirExists(baseView) || !dirExists(cmpView) {
			continue
		}
		diffView := ""
		if diffDir != "" {
			diffView = filepath.Join(diffDir, view)
			os.MkdirAll(diffView, 0755)
		}
		fmt.Printf("\n=== %s ===\n", view)
		changed := compareDir(baseView, cmpView, diffView)
		total += changed
	}

	fmt.Println()
	if total > 0 {
		fmt.Printf("FAIL: %d screenshots differ\n", total)
		os.Exit(1)
	}
	fmt.Println("PASS: All screenshots match")
}

func cmdCIDiff(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: visual-diff ci-diff <local-dir> [skip-build]")
		os.Exit(1)
	}
	localDir := args[0]
	skipBuild := ""
	if len(args) > 1 {
		skipBuild = args[1]
	}

	fmt.Println("Finding latest stable build...")
	build := findLatestBuild("stable")
	if build == "" {
		fmt.Println("No stable build found, skipping visual diff")
		return
	}

	if skipBuild != "" && skipBuild == build {
		fmt.Printf("Skipping visual diff against stable build #%s\n", build)
		return
	}

	fmt.Printf("Stable build: #%s\n", build)

	stableDir := "/tmp/visual-diff-stable"
	os.RemoveAll(stableDir)

	for _, view := range views {
		url := fmt.Sprintf("%s/%s-%s/distro/%s/screenshot", ciFiles, build, arch, view)
		dir := filepath.Join(stableDir, view)
		os.MkdirAll(dir, 0755)
		fmt.Printf("Downloading stable %s screenshots...\n", view)
		downloadScreenshots(url, dir)
	}

	fmt.Println()
	total := 0
	for _, view := range views {
		baseView := filepath.Join(stableDir, view)
		cmpView := filepath.Join(localDir, view)
		if !dirExists(baseView) || !dirExists(cmpView) {
			continue
		}
		fmt.Printf("\n=== %s ===\n", view)
		changed := compareDir(baseView, cmpView, "")
		total += changed
	}

	fmt.Println()
	if total > 0 {
		fmt.Printf("FAIL: %d screenshots differ\n", total)
		os.Exit(1)
	}
	fmt.Println("PASS: All screenshots match stable")
}

func cmdDownload(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: visual-diff download <branch> <output-dir>")
		os.Exit(1)
	}
	branch := args[0]
	outputDir := args[1]

	fmt.Printf("Finding latest build for '%s'...\n", branch)
	build := findLatestBuild(branch)
	if build == "" {
		fmt.Fprintf(os.Stderr, "ERROR: No build found for branch '%s'\n", branch)
		os.Exit(1)
	}

	fmt.Printf("Branch: %s (build #%s)\n", branch, build)

	os.RemoveAll(outputDir)
	for _, view := range views {
		url := fmt.Sprintf("%s/%s-%s/distro/%s/screenshot", ciFiles, build, arch, view)
		dir := filepath.Join(outputDir, view)
		os.MkdirAll(dir, 0755)
		fmt.Printf("\n=== %s screenshots ===\n", view)
		downloadScreenshots(url, dir)
	}
}

func findLatestBuild(branch string) string {
	url := fmt.Sprintf("%s/builds?limit=5&branch=%s", ciAPI, branch)
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var builds []Build
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

func listRemoteFiles(url string) []string {
	resp, err := http.Get(url + "/")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var files []FileEntry
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil
	}

	var names []string
	for _, f := range files {
		if strings.HasSuffix(f.Name, ".png") {
			names = append(names, f.Name)
		}
	}
	return names
}

func downloadScreenshots(url, dir string) {
	files := listRemoteFiles(url)
	for _, name := range files {
		target := filepath.Join(dir, name)
		if _, err := os.Stat(target); err == nil {
			continue
		}
		downloadFile(url+"/"+name, target)
		fmt.Printf("  %s\n", name)
	}
}

func downloadFile(url, path string) error {
	resp, err := http.Get(url)
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

func shouldSkip(name string) bool {
	if strings.HasPrefix(name, "exception") {
		return true
	}
	if strings.Contains(name, "_unstable") {
		return true
	}
	return false
}

func compareDir(baseDir, cmpDir, diffDir string) int {
	changed := 0
	identical := 0
	missing := 0

	entries, _ := os.ReadDir(baseDir)
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".png") {
			continue
		}
		if shouldSkip(name) {
			continue
		}

		basePath := filepath.Join(baseDir, name)
		cmpPath := filepath.Join(cmpDir, name)

		if _, err := os.Stat(cmpPath); os.IsNotExist(err) {
			fmt.Printf("  MISSING: %s\n", name)
			missing++
			continue
		}

		pixels, err := compareImages(basePath, cmpPath)
		if err != nil {
			fmt.Printf("  ERROR: %s (%v)\n", name, err)
			changed++
			continue
		}

		if pixels <= pixelThreshold {
			identical++
		} else {
			fmt.Printf("  CHANGED: %s (%d pixels differ)\n", name, pixels)
			changed++
		}
	}

	// Check for new files
	cmpEntries, _ := os.ReadDir(cmpDir)
	for _, e := range cmpEntries {
		name := e.Name()
		if !strings.HasSuffix(name, ".png") || shouldSkip(name) {
			continue
		}
		basePath := filepath.Join(baseDir, name)
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			fmt.Printf("  NEW: %s\n", name)
			missing++
		}
	}

	fmt.Printf("  Identical: %d, Changed: %d, Missing/New: %d\n", identical, changed, missing)
	return changed
}

func compareImages(path1, path2 string) (int, error) {
	img1, err := loadPNG(path1)
	if err != nil {
		return 0, fmt.Errorf("loading %s: %w", path1, err)
	}
	img2, err := loadPNG(path2)
	if err != nil {
		return 0, fmt.Errorf("loading %s: %w", path2, err)
	}

	b1 := img1.Bounds()
	b2 := img2.Bounds()

	if b1.Dx() != b2.Dx() || b1.Dy() != b2.Dy() {
		return b1.Dx() * b1.Dy(), nil
	}

	fuzz := fuzzPercent / 100.0 * 65535.0
	diffCount := 0

	for y := b1.Min.Y; y < b1.Max.Y; y++ {
		for x := b1.Min.X; x < b1.Max.X; x++ {
			r1, g1, b1c, a1 := img1.At(x, y).RGBA()
			r2, g2, b2c, a2 := img2.At(x, y).RGBA()

			if colorDiff(r1, r2) > fuzz ||
				colorDiff(g1, g2) > fuzz ||
				colorDiff(b1c, b2c) > fuzz ||
				colorDiff(a1, a2) > fuzz {
				diffCount++
			}
		}
	}

	return diffCount, nil
}

func colorDiff(a, b uint32) float64 {
	return math.Abs(float64(a) - float64(b))
}

func loadPNG(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
