package dto

type UserLogin struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=1,max=64"`
}
