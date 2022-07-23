package dto

type CreateGym struct {
	Title   string `json:"title" binding:"required"`
	Address string `json:"address" binding:"required"`
}
