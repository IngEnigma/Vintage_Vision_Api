package request

type CreateProfileRequest struct {
	Name      string `json:"name" binding:"required"`
	AvatarUrl string `json:"avatar_url" binding:"required,url"`
}

type UpdateProfileRequest struct {
	Name      *string `json:"name" binding:"omitempty,min=1,max=10"`
	AvatarURL *string `json:"avatar_url" binding:"omitempty,url"`
}
