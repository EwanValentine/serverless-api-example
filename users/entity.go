package users

// User -
type User struct {
	ID    string `json:"id"`
	Email string `json:"email" validate:"email,required"`
	Name  string `json:"name" validate:"required,gte=1,lte=50"`
	Age   uint32 `json:"age" validate:"required,gte=0,lte=130"`
}

// UpdateUser -
type UpdateUser struct {
	Email string `json:"email"`
	Name  string `json:"name" validate:"gte=1,lte=50"`
	Age   uint32 `json:"age" validate:"gte=0,lte=130"`
}
