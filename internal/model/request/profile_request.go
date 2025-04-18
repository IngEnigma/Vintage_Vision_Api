package request

type CreateProfileRequest struct {
	Name      string `json:"name" binding:"required"`
	AvatarUrl string `json:"avatar_url"`
}

type UpdateProfileRequest struct {
	Name      *string `json:"name" binding:"omitempty,min=1"`
	AvatarURL *string `json:"avatar" binding:"omitempty,url"`
}
