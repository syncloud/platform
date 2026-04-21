package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type App struct {
	ci         *CIClient
	comparator *Comparator
	pictures   *Pictures
}

func (a *App) Diff(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: visual-diff diff <base-dir> <compare-dir> [diff-dir]")
		os.Exit(1)
	}
	baseDir, cmpDir := args[0], args[1]
	diffDir := ""
	if len(args) > 2 {
		diffDir = args[2]
	}
	changed, compared := a.comparator.CompareTrees(baseDir, cmpDir, diffDir)
	a.reportResult(changed, compared, "")
}

func (a *App) CIDiff(args []string) {
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
	build := a.ci.FindLatestBuild("stable")
	if build == "" {
		fmt.Println("No stable build found, skipping visual diff")
		return
	}
	if skipBuild != "" && skipBuild == build {
		fmt.Printf("Skipping visual diff against stable build #%s\n", build)
		return
	}
	fmt.Printf("Stable build: #%s\n", build)

	stableDir := filepath.Join(os.TempDir(), "visual-diff-stable")
	os.RemoveAll(stableDir)
	if err := a.ci.DownloadBuild(build, stableDir); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR downloading stable: %v\n", err)
		os.Exit(1)
	}

	changed, compared := a.comparator.CompareTrees(stableDir, localDir, "")
	a.reportResult(changed, compared, "stable")
}

func (a *App) Download(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: visual-diff download <branch> [output-dir]")
		os.Exit(1)
	}
	branch := args[0]
	outDir := ""
	if len(args) > 1 {
		outDir = args[1]
	} else {
		outDir = filepath.Join(os.TempDir(), "visual-diff-cache", branch)
	}
	if err := a.ci.DownloadBranch(branch, outDir); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	a.pictures.CopyBranch(outDir, branch, a.comparator.views)
}

func (a *App) Branches(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: visual-diff branches <base-branch> <compare-branch>")
		os.Exit(1)
	}
	base, cmp := args[0], args[1]

	cache := filepath.Join(os.TempDir(), "visual-diff-cache")
	baseDir := filepath.Join(cache, base)
	cmpDir := filepath.Join(cache, cmp)
	diffDir := filepath.Join(cache, "diff")
	os.RemoveAll(diffDir)

	if err := a.ci.DownloadBranch(base, baseDir); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR downloading %s: %v\n", base, err)
		os.Exit(1)
	}
	if err := a.ci.DownloadBranch(cmp, cmpDir); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR downloading %s: %v\n", cmp, err)
		os.Exit(1)
	}

	changed, compared := a.comparator.CompareTrees(baseDir, cmpDir, diffDir)
	a.pictures.Copy(baseDir, cmpDir, diffDir, base, cmp, a.comparator.views)
	a.reportResult(changed, compared, "")
}

func (a *App) reportResult(changed, compared int, label string) {
	suffix := ""
	if label != "" {
		suffix = " " + label
	}
	fmt.Println()
	if compared == 0 {
		fmt.Printf("FAIL: no screenshots were compared%s\n", suffix)
		os.Exit(2)
	}
	if changed > 0 {
		fmt.Printf("FAIL: %d screenshots differ%s\n", changed, suffix)
		os.Exit(1)
	}
	fmt.Printf("PASS: %d screenshots match%s\n", compared, suffix)
}
