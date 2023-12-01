package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChange_ParseApp(t *testing.T) {
	assert.Equal(t, "matrix", ParseApp(`Refresh \"matrix\" snap from \"latest/stable\" channel`))
	assert.Equal(t, "matrix", ParseApp(`"Install \"matrix\" snap from \"latest/stable\" channel`))
	assert.Equal(t, "matrix", ParseApp(`"Remove \"matrix\" snap from \"latest/stable\" channel`))
	assert.Equal(t, "unknown", ParseApp(`"doing something`))
}

func TestChange_ParseAction(t *testing.T) {
	assert.Equal(t, "Upgrading", ParseAction(`Refresh \"matrix\" snap from \"latest/stable\" channel`))
	assert.Equal(t, "Installing", ParseAction(`Install \"matrix\" snap from \"latest/stable\" channel`))
	assert.Equal(t, "Removing", ParseAction(`Remove \"matrix\" snap from \"latest/stable\" channel`))
	assert.Equal(t, "Unknown", ParseAction(`doing something`))
}
