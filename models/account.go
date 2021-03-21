package models

type AccountCheckoutResp struct {
	Code    int     `json:"code"`
	Message string  `json:"msg"`
	Data    Account `json:"data"`
}

type Account struct {
	Username string
	Password string
}
