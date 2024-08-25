package apimodels

type RegisterReq struct {
	Nickname string `json:"nickname" validate:"required,min=1,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokensResp struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type VerifyEmailReq struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
