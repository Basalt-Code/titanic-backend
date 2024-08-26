package apimodels

type RegisterReq struct {
	Nickname string `json:"nickname" validate:"required,min=1,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}
