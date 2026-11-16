package rest

type UserPasskeyRemoveRequest struct {
	Username string `json:"username"`
	ID       int64  `json:"id"`
}
