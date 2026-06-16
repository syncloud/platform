package rest

type GroupSetMemberRequest struct {
	Group    string `json:"group"`
	Username string `json:"username"`
	Member   bool   `json:"member"`
}
