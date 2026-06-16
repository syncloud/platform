package rest

type GroupMemberRequest struct {
	Group    string `json:"group"`
	Username string `json:"username"`
}
