package dto

type (
	CreateUserRequest struct {
		Email    string `json:"email" valid:"required,email"`
		Password string `json:"password" valid:"required"`
		Role     string `json:"role" valid:"required,in(ADMIN|USER)"`

		CreatedBy string `json:"-"`
	}
)
