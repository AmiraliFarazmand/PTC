package ports

import "github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"

type ZeebeProcessManager interface {
	StartSignupProcess(username, password string) error
	StartLoginProcess(username, password string) (*domain.ProcessVariables, error)
	StartPurchaseProcess(userID string, amount int, address string) (*domain.PurchaseProcessVariables, error)
}
