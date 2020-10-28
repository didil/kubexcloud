package responses

// LoginUser response
type LoginUser struct {
	Token string `json:"token"`
}

// ListUser response
type ListUser struct {
	Users []ListUserEntry `json:"users"`
}

type ListUserEntry struct {
	Name string `json:"name"`
	Role string `json:"role"`
}
