package response

import "github.com/olivercruznaguit/inventory-system/internal/model"

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

func NewUserResponse(user model.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
}
