package auth

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func writeNotification(t *testing.T, content string) string {
	file := path.Join(t.TempDir(), "notification.txt")
	assert.NoError(t, os.WriteFile(file, []byte(content), 0600))
	return file
}

func TestElevationCode_Read(t *testing.T) {
	file := writeNotification(t, "Date: 2026-11-14 08:48:09 +0000 UTC\nRecipient: \"user\" <user@example.com>\nSubject: Confirm your identity\nABCD1234\n\n")
	code, err := NewElevationCode(file).Read("user@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "ABCD1234", code)
}

func TestElevationCode_Read_AutheliaRecipient(t *testing.T) {
	file := writeNotification(t, "Date: 2026-11-15 04:05:28.123 +0900 KST\nRecipient: {user user@buster-amd64.redirect}\nSubject: Confirm your identity\nABCD1234\n\n")
	code, err := NewElevationCode(file).Read("user@buster-amd64.redirect")
	assert.NoError(t, err)
	assert.Equal(t, "ABCD1234", code)
}

func TestElevationCode_Read_AutheliaOtherRecipient(t *testing.T) {
	file := writeNotification(t, "Date: 2026-11-15 04:05:28.123 +0900 KST\nRecipient: {auser auser@buster-amd64.redirect}\nSubject: Confirm your identity\nABCD1234\n\n")
	_, err := NewElevationCode(file).Read("user@buster-amd64.redirect")
	assert.Error(t, err)
}

func TestElevationCode_Read_OtherRecipient(t *testing.T) {
	file := writeNotification(t, "Date: 2026-11-14 08:48:09 +0000 UTC\nRecipient: \"other\" <other@example.com>\nSubject: Confirm your identity\nABCD1234\n\n")
	_, err := NewElevationCode(file).Read("user@example.com")
	assert.Error(t, err)
}

func TestElevationCode_Read_RecipientSuffix(t *testing.T) {
	file := writeNotification(t, "Date: 2026-11-14 08:48:09 +0000 UTC\nRecipient: \"auser\" <auser@example.com>\nSubject: Confirm your identity\nABCD1234\n\n")
	_, err := NewElevationCode(file).Read("user@example.com")
	assert.Error(t, err)
}

func TestElevationCode_Read_NoCode(t *testing.T) {
	file := writeNotification(t, "Date: 2026-11-14 08:48:09 +0000 UTC\nRecipient: \"user\" <user@example.com>\nSubject: Confirm your identity\n\n")
	_, err := NewElevationCode(file).Read("user@example.com")
	assert.Error(t, err)
}

func TestElevationCode_Read_Missing(t *testing.T) {
	_, err := NewElevationCode(path.Join(t.TempDir(), "missing.txt")).Read("user@example.com")
	assert.Error(t, err)
}
