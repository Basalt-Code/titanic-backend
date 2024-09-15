package apimodels

type RegisterReq struct {
	Username string `json:"username" validate:"required,min=1,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	Role     string `json:"role" validate:"required,role"`
}

type LoginReq struct {
	Username string `json:"username" validate:"required,username"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type TokensResp struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
