package auth

import (
	"fmt"
	"os"
	"strings"
)

const (
	recipientHeader = "Recipient: "
	subjectHeader   = "Subject: "
)

type ElevationCode struct {
	notificationFile string
}

func NewElevationCode(notificationFile string) *ElevationCode {
	return &ElevationCode{notificationFile: notificationFile}
}

func (e *ElevationCode) Read(email string) (string, error) {
	content, err := os.ReadFile(e.notificationFile)
	if err != nil {
		return "", fmt.Errorf("read notification: %w", err)
	}
	recipient := ""
	subject := false
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, recipientHeader):
			recipient = strings.TrimPrefix(line, recipientHeader)
		case strings.HasPrefix(line, subjectHeader):
			subject = true
		case subject && line != "":
			if !addressedTo(recipient, email) {
				return "", fmt.Errorf("notification is not for %s", email)
			}
			return line, nil
		}
	}
	return "", fmt.Errorf("no code in notification")
}

func addressedTo(recipient, email string) bool {
	for _, field := range strings.Fields(recipient) {
		if strings.Trim(field, "{}<>\"") == email {
			return true
		}
	}
	return false
}
