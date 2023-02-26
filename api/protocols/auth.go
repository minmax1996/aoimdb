package protocols

type AuthRequest struct {
	User     *string `json:"user"`
	Password *string `json:"password"`
	Token    *string `json:"token"`
}

type AuthResponse struct {
	Message string `json:"message"`
}
