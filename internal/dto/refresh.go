package dto

type RefreshOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token"`
}
