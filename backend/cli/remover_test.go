package cli

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestRemove(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(path.Join(dir, "1.log"), []byte(""), 644)
	_ = os.WriteFile(path.Join(dir, "2.txt"), []byte(""), 644)
	_ = os.WriteFile(path.Join(dir, "3.log"), []byte(""), 644)

	err := Remove(fmt.Sprintf("%s/*.log", dir))
	assert.NoError(t, err)

	entries, err := os.ReadDir(dir)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(entries))
	assert.Equal(t, "2.txt", entries[0].Name())

}
