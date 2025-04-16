package ports


type PurchaseService interface {
    CreatePurchase(userID string, amount int, address string) error
    ConfirmPayment(purchaseID string) error
}