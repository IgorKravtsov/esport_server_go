package dto

type UserRegister struct {
	Name     string `json:"name" binding:"required,min=1,max=64"`
	Email    string `json:"email" binding:"required,email,max=63"`
	Password string `json:"password" binding:"required,min=1,max=64"`
}
