package backend

type NewEmailRequest struct {
	Email string
}

type RegisterRequest struct {
	Username string
	Password string
}

type LoginRequest struct {
	*RegisterRequest
}
