package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type Comparator struct {
	pixelThreshold int
	fuzz           float64
	views          []string
}

func NewComparator(pixelThreshold int, fuzzPercent float64, views []string) *Comparator {
	return &Comparator{
		pixelThreshold: pixelThreshold,
		fuzz:           fuzzPercent / 100.0 * 65535.0,
		views:          views,
	}
}

func (c *Comparator) CompareTrees(baseRoot, cmpRoot, diffRoot string) (int, int) {
	totalChanged, totalCompared := 0, 0
	for _, view := range c.views {
		baseView := resolveView(baseRoot, view)
		cmpView := resolveView(cmpRoot, view)
		if baseView == "" || cmpView == "" {
			fmt.Printf("\n=== %s ===\n  SKIP: base=%s compare=%s\n", view, presence(baseView), presence(cmpView))
			continue
		}
		diffView := ""
		if diffRoot != "" {
			diffView = filepath.Join(diffRoot, view)
			os.MkdirAll(diffView, 0755)
		}
		fmt.Printf("\n=== %s ===\n", view)
		changed, compared := c.compareDir(baseView, cmpView, diffView)
		totalChanged += changed
		totalCompared += compared
	}
	return totalChanged, totalCompared
}

func (c *Comparator) compareDir(baseDir, cmpDir, diffDir string) (int, int) {
	changed, identical, missing, compared := 0, 0, 0, 0

	entries, _ := os.ReadDir(baseDir)
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".png") || shouldSkip(name) {
			continue
		}
		cmpPath := filepath.Join(cmpDir, name)
		if _, err := os.Stat(cmpPath); os.IsNotExist(err) {
			fmt.Printf("  MISSING: %s\n", name)
			missing++
			continue
		}
		compared++
		diffPath := ""
		if diffDir != "" {
			diffPath = filepath.Join(diffDir, name)
		}
		pixels, err := c.compareImages(filepath.Join(baseDir, name), cmpPath, diffPath)
		if err != nil {
			fmt.Printf("  ERROR: %s (%v)\n", name, err)
			changed++
			continue
		}
		if pixels <= c.pixelThreshold {
			identical++
			if diffPath != "" {
				os.Remove(diffPath)
			}
		} else {
			fmt.Printf("  CHANGED: %s (%d pixels differ)\n", name, pixels)
			changed++
		}
	}

	cmpEntries, _ := os.ReadDir(cmpDir)
	for _, e := range cmpEntries {
		name := e.Name()
		if !strings.HasSuffix(name, ".png") || shouldSkip(name) {
			continue
		}
		if _, err := os.Stat(filepath.Join(baseDir, name)); os.IsNotExist(err) {
			fmt.Printf("  NEW: %s\n", name)
			missing++
		}
	}

	fmt.Printf("  Identical: %d, Changed: %d, Missing/New: %d\n", identical, changed, missing)
	return changed, compared
}

func (c *Comparator) compareImages(path1, path2, diffPath string) (int, error) {
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
		if diffPath != "" {
			writeDiffPNG(img2, nil, diffPath)
		}
		return b1.Dx() * b1.Dy(), nil
	}

	var mask []bool
	if diffPath != "" {
		mask = make([]bool, b1.Dx()*b1.Dy())
	}
	diffCount := 0
	for y := b1.Min.Y; y < b1.Max.Y; y++ {
		for x := b1.Min.X; x < b1.Max.X; x++ {
			r1, g1, b1c, a1 := img1.At(x, y).RGBA()
			r2, g2, b2c, a2 := img2.At(x, y).RGBA()
			if colorDiff(r1, r2) > c.fuzz ||
				colorDiff(g1, g2) > c.fuzz ||
				colorDiff(b1c, b2c) > c.fuzz ||
				colorDiff(a1, a2) > c.fuzz {
				diffCount++
				if mask != nil {
					mask[(y-b1.Min.Y)*b1.Dx()+(x-b1.Min.X)] = true
				}
			}
		}
	}
	if diffPath != "" && diffCount > 0 {
		writeDiffPNG(img1, mask, diffPath)
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

func writeDiffPNG(base image.Image, mask []bool, path string) error {
	b := base.Bounds()
	out := image.NewRGBA(b)
	red := color.RGBA{R: 255, A: 255}
	w := b.Dx()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			idx := (y-b.Min.Y)*w + (x - b.Min.X)
			if mask != nil && idx < len(mask) && mask[idx] {
				out.Set(x, y, red)
				continue
			}
			r, g, bl, a := base.At(x, y).RGBA()
			gray := uint8((r + g + bl) / 3 >> 8)
			out.Set(x, y, color.RGBA{R: gray, G: gray, B: gray, A: uint8(a >> 8)})
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, out)
}

func resolveView(dir, view string) string {
	if dir == "" {
		return ""
	}
	for _, candidate := range []string{
		filepath.Join(dir, view, "screenshot"),
		filepath.Join(dir, view),
	} {
		if hasPNGs(candidate) {
			return candidate
		}
	}
	return ""
}

func hasPNGs(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".png") {
			return true
		}
	}
	return false
}

func presence(path string) string {
	if path == "" {
		return "missing"
	}
	return path
}

func shouldSkip(name string) bool {
	return strings.HasPrefix(name, "exception") || strings.Contains(name, "_unstable")
}
