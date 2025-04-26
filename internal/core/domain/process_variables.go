package domain

type ProcessVariables struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	IsValid    bool   `json:"isValid"`
	LoginValid bool   `json:"loginValid"`
	Token      string `json:"token"`
	Error      string `json:"error"`
}
