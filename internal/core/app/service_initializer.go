package app

import (
	"github.com/AmiraliFarazmand/PTC_Task/internal/ports"
)

func InitializeUserService(userRepo ports.UserRepository) ports.UserService {
	return &UserServiceImpl{
		UserRepo: userRepo,
	}
}

func InitializePurchaseService(purchaseRepo ports.PurchaseRepository) ports.PurchaseService {
	return &PurchaseServiceImpl{
		PurchaseRepo: purchaseRepo,
	}
}
