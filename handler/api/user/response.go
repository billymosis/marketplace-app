package user

type loginUserResponse struct {
	Message string `json:"message"`
	Data    struct {
		Username    string `json:"username"`
		Name        string `json:"name"`
		AccessToken string `json:"accessToken"`
	} `json:"data"`
}

type createUserResponse struct {
	Message string `json:"message"`
	Data    struct {
		Username    string `json:"username"`
		Name        string `json:"name"`
		AccessToken string `json:"accessToken"`
	} `json:"data"`
}
