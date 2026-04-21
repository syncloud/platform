package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Pictures struct {
	root string
}

func NewPictures(root string) *Pictures {
	return &Pictures{root: root}
}

func (p *Pictures) Copy(baseDir, cmpDir, diffDir, baseLabel, cmpLabel string, views []string) {
	if p.root == "" {
		return
	}
	target := filepath.Join(p.root, "screenshot-diff")
	os.RemoveAll(target)

	sets := []struct {
		dst, src string
	}{
		{filepath.Join(target, "base"), baseDir},
		{filepath.Join(target, "compare"), cmpDir},
		{filepath.Join(target, "diff"), diffDir},
	}
	for _, s := range sets {
		p.copyViews(s.src, s.dst, views)
	}
	fmt.Printf("\nPictures/screenshot-diff/ (base=%s, compare=%s)\n", baseLabel, cmpLabel)
}

func (p *Pictures) CopyBranch(srcDir, label string, views []string) {
	if p.root == "" {
		return
	}
	target := filepath.Join(p.root, "syncloud-"+label)
	os.RemoveAll(target)
	p.copyViews(srcDir, target, views)
	fmt.Printf("\nPictures/syncloud-%s/\n", label)
}

func (p *Pictures) copyViews(src, dst string, views []string) {
	for _, view := range views {
		srcView := resolveView(src, view)
		if srcView == "" {
			continue
		}
		if err := os.MkdirAll(dst, 0755); err != nil {
			continue
		}
		entries, _ := os.ReadDir(srcView)
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".png") {
				continue
			}
			copyFile(filepath.Join(srcView, e.Name()), filepath.Join(dst, e.Name()))
		}
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
