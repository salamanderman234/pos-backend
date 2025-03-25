package forms

type FormLogin struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type FormResendTwoFactor struct {
	Seed string `json:"seed" form:"seed" validate:"required"`
}

type FormVerifyTwoFactor struct {
	Seed string `json:"seed" form:"seed" validate:"required"`
	Code string `json:"code" form:"code" validate:"required"`
}

type FormVerifyUser struct {
	Username string `json:"username" form:"username" validate:"required"`
	Key      string `json:"key" form:"key" validate:"required"`
}

type FormResetPassword struct {
	Username    string `json:"username" form:"username" validate:"required"`
	Key         string `json:"key" form:"key" validate:"required"`
	NewPassword string `json:"new_password" form:"new_password" validate:"required,password"`
}

type FormResendEmail struct {
	Username string `json:"username" form:"username" validate:"required"`
}
