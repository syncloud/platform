package backup

import (
	"io/ioutil"
		"os"
	"log"
	"path/filepath"

	"testing"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	tmpfn := filepath.Join(dir, "tmpfile")
	if err := ioutil.WriteFile(tmpfn, []byte(""), 0666); err != nil {
		log.Fatal(err)
	}
	list, err := List(dir)
	assert.Nil(t, err)
	assert.Equal(t, list, []string{"tmpfile"}, "")
} 