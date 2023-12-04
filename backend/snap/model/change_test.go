package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseApp(t *testing.T) {
	assert.Equal(t, "matrix", ParseApp(`Refresh "matrix" snap from "latest/stable" channel`))
	assert.Equal(t, "matrix", ParseApp(`"Install "matrix" snap from "latest/stable" channel`))
	assert.Equal(t, "matrix", ParseApp(`"Remove "matrix" snap from "latest/stable" channel`))
	assert.Equal(t, "unknown", ParseApp(`"doing something`))
}

func TestParseAction(t *testing.T) {
	assert.Equal(t, "Upgrading", ParseAction(`Refresh "matrix" snap from "latest/stable" channel`))
	assert.Equal(t, "Installing", ParseAction(`Install "matrix" snap from "latest/stable" channel`))
	assert.Equal(t, "Removing", ParseAction(`Remove "matrix" snap from "latest/stable" channel`))
	assert.Equal(t, "Unknown", ParseAction(`doing something`))
}

func TestCalcPercentage(t *testing.T) {
	assert.Equal(t, int64(0), CalcPercentage(0, 0))
	assert.Equal(t, int64(0), CalcPercentage(1, 1))
	assert.Equal(t, int64(9), CalcPercentage(12, 123))
}
