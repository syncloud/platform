package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	ci := NewCIClient(
		"http://ci.syncloud.org:8080/api/repos/syncloud/platform",
		"http://ci.syncloud.org:8081/files/platform",
		"amd64",
		&http.Client{Timeout: 60 * time.Second},
	)
	comparator := NewComparator(100, 5.0, []string{"desktop", "mobile"})
	pictures := NewPictures(picturesRoot())

	app := &App{ci: ci, comparator: comparator, pictures: pictures}

	switch os.Args[1] {
	case "diff":
		app.Diff(os.Args[2:])
	case "ci-diff":
		app.CIDiff(os.Args[2:])
	case "download":
		app.Download(os.Args[2:])
	case "branches":
		app.Branches(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Usage: visual-diff <command> [args]")
	fmt.Println("Commands:")
	fmt.Println("  diff <base-dir> <compare-dir> [diff-dir]    Compare two local screenshot dirs")
	fmt.Println("  ci-diff <local-dir> [skip-build]            Compare local artifact/distro against stable")
	fmt.Println("  download <branch> <output-dir>              Download screenshots from CI")
	fmt.Println("  branches <base-branch> <compare-branch>     Download and compare two branches")
}

func picturesRoot() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	path := filepath.Join(home, "storage", "pictures")
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return path
	}
	return ""
}
