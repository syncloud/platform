package auth

import "time"

type Passkey struct {
	ID          int64      `json:"id"`
	Description string     `json:"description"`
	Created     time.Time  `json:"created"`
	LastUsed    *time.Time `json:"last_used"`
}
