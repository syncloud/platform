// Package lcd provides LCD display support for devices with SSD1306 OLED displays
//
// Based on https://github.com/ChandlerSwift/odroidhc4-display (MIT License)
// Copyright (c) 2021 Chandler Swift
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
package lcd

import (
	"fmt"
	"image"
	"log"
	"net"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ssd1306"
	"periph.io/x/devices/v3/ssd1306/image1bit"
	"periph.io/x/host/v3"
)

// Display manages the LCD display output
type Display struct {
	updateInterval time.Duration
}

// NewDisplay creates a new LCD display manager
func NewDisplay() *Display {
	return &Display{
		updateInterval: time.Second,
	}
}

// Start starts the LCD display service in a goroutine (non-blocking)
// Implements the Service interface
func (d *Display) Start() error {
	// Check if I2C device is available before starting
	if !d.isAvailable() {
		log.Println("LCD I2C device not found, skipping LCD display")
		return nil
	}

	go func() {
		log.Println("Starting LCD display service...")
		if err := d.run(); err != nil {
			log.Printf("LCD display error: %v", err)
		}
	}()
	return nil
}

// isAvailable checks if the I2C device is available
func (d *Display) isAvailable() bool {
	// Try to initialize periph and open I2C bus
	if _, err := host.Init(); err != nil {
		return false
	}
	b, err := i2creg.Open("")
	if err != nil {
		return false
	}
	_ = b.Close()
	return true
}

// run starts the LCD display loop (blocking)
func (d *Display) run() error {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		return fmt.Errorf("failed to initialize periph: %w", err)
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	b, err := i2creg.Open("")
	if err != nil {
		return fmt.Errorf("failed to open I2C bus: %w", err)
	}
	defer b.Close()

	dev, err := ssd1306.NewI2C(b, &ssd1306.Opts{
		W:             128,
		H:             64,
		Rotated:       true,
		Sequential:    false,
		SwapTopBottom: false,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize ssd1306: %w", err)
	}

	f := basicfont.Face7x13
	drawer := font.Drawer{
		Src:  &image.Uniform{image1bit.On},
		Face: f,
		Dot:  fixed.P(0, f.Height),
	}

	ticker := time.NewTicker(d.updateInterval)
	defer ticker.Stop()

	for {
		t := <-ticker.C

		lines := []string{
			t.Format("Jan 2 3:04:05 PM"),
			getIPAddrString(),
			getCPUString(),
			getMemString(),
			getHDDString(),
		}

		img := image1bit.NewVerticalLSB(dev.Bounds()) // reset canvas per frame
		drawer.Dst = img
		for i, s := range lines {
			drawer.Dot = fixed.P(0, (f.Height-1)*(i+1))
			drawer.DrawString(s)
		}

		if err := dev.Draw(dev.Bounds(), img, image.Point{}); err != nil {
			log.Printf("error drawing to display: %v", err)
		}
	}
}

// formatSize formats bytes to human-readable format
// based on https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func formatSize(size uint64, unit uint64) string {
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}
	div, suffix := unit, 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		suffix++
	}
	value := float64(size) / float64(div)
	var fmtStr string
	if value >= 100 {
		fmtStr = "%.f%c"
	} else {
		fmtStr = "%.1f%c"
	}

	return fmt.Sprintf(fmtStr, value, "kMGTP"[suffix])
}

func getMemString() string {
	v, err := mem.VirtualMemory()
	if err != nil {
		return "MEM: unavailable"
	}
	return fmt.Sprintf("MEM: %v/%v", formatSize(v.Used, 1024), formatSize(v.Total, 1024))
}

func getCPUString() string {
	v, err := cpu.Percent(0, false)
	if err != nil {
		return "CPU: unavailable"
	}
	l, err := load.Avg()
	if err != nil {
		return fmt.Sprintf("CPU: %.f%%", v[0])
	}
	// Unfortunately, the screen just isn't wide enough to include Load15
	return fmt.Sprintf("CPU: %.f%% (%.1f %.1f)", v[0], l.Load1, l.Load5)
}

// getHDDString returns data about the biggest mounted partition.
func getHDDString() string {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return "DISK: unavailable"
	}
	biggestDiskSize := uint64(0)
	biggestDiskUsed := uint64(0)
	biggestDiskName := ""
	for _, partition := range partitions {
		d, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}
		if d.Total > biggestDiskSize {
			biggestDiskName = partition.Mountpoint
			biggestDiskUsed = d.Used
			biggestDiskSize = d.Total
		}
	}
	if biggestDiskName == "" {
		return "DISK: unavailable"
	}
	return fmt.Sprintf("%v: %v/%v", biggestDiskName, formatSize(biggestDiskUsed, 1000), formatSize(biggestDiskSize, 1000))
}

func getIPAddrString() string {
	// https://stackoverflow.com/a/37382208/3814663
	// Note that since this is UDP, no connection is actually established.
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "IP: Network down"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return fmt.Sprintf("IP: %v", localAddr.IP)
}
