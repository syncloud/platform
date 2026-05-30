package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToSyncloudApp_IconUrl(t *testing.T) {
	snap := &Snap{Name: "Test", Summary: "Summary", Channel: "stable", Version: "1", Type: "app", Apps: nil}
	app := snap.toSyncloudApp("url")
	assert.Equal(t, "/rest/proxy/image?channel=stable&app=Test", app.App.Icon)
}

func TestToSyncloudApp_NormalizeIconUrl_AfterLocalAmendInstall(t *testing.T) {
	snap := &Snap{Name: "Test", Summary: "Summary", Channel: "latest/stable", Version: "1", Type: "app", Apps: nil}
	app := snap.toSyncloudApp("url")
	assert.Equal(t, "/rest/proxy/image?channel=stable&app=Test", app.App.Icon)
}

func TestSnap_IconUrl(t *testing.T) {
	snap := &Snap{Media: []SnapMedia{
		{Type: "screenshot", Url: "http://store/shot.png"},
		{Type: "icon", Url: "http://store/v2/apps/stable/test/icon.png"},
	}}
	assert.Equal(t, "http://store/v2/apps/stable/test/icon.png", snap.IconUrl())
}

func TestSnap_IconUrl_None(t *testing.T) {
	snap := &Snap{}
	assert.Equal(t, "", snap.IconUrl())
}
