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
