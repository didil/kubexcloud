package requests

// LoginUser request
type LoginUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// CreateUser request
type CreateUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
