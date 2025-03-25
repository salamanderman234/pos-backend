package forms

type FormLogin struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type FormResendTwoFactor struct {
	Seed string `json:"seed" form:"seed" validate:"required"`
}
