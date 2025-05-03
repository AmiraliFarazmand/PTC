package ports

type PurchaseService interface {
	CreatePurchase(userID string, amount int, address string) (string, error)
	ConfirmPayment(purchaseID string, userID string) error
	CancelUnpaidPurchase(purchaseID string)	error
}
