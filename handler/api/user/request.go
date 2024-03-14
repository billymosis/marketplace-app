package user

type loginUserRequest struct {
	Username string `json:"username" validate:"required,min=5,max=15"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}

type createUserRequest struct {
	Username string `json:"username" validate:"required,min=5,max=15"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}
