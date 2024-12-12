package dto

// UserRegisterRequest defines the structure of the data required for user registration.
// It includes validation tags to ensure the presence of mandatory fields and a valid email format.
type UserRegisterRequest struct {
	Nickname string `json:"nickname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

// UserResponse defines the structure of the response returned during user operations.
// It usually includes messages to inform the client of the operation's success or failure.
type UserResponse struct {
	Message string `json:"message"`
}
