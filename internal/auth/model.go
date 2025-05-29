package auth

type RegisterInput struct {
	Email string `json:"email"`
	Password string `json:"password"`
	FullName *string `json:"full_name,omitempty"`
}
