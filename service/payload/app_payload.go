package payload

type LoginRequestBodyValue struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequestBody struct {
	Request LoginRequestBodyValue `json:"request"`
}
