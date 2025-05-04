package domain

type AuthProcessVariables struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	IsValid    bool   `json:"isValid"`
	Token      string `json:"token"`
	Error      string `json:"error"`
}

type PurchaseProcessVariables struct {
	UserID     string `json:"user_id"`
	Amount     int    `json:"amount"`
	Address    string `json:"address"`
	PurchaseID string `json:"purchase_id"`
	IsValid    bool   `json:"isValid"`
	Error      string `json:"error,omitempty"`
}
