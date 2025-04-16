package app

import "github.com/AmiraliFarazmand/PTC_Task/internal/domain"

func InitializePurchaseService(purchaseRepo domain.PurchaseRepository) PurchaseServiceImpl {
    return PurchaseServiceImpl{PurchaseRepo: purchaseRepo}
}

func InitializeUserService(userRepo domain.UserRepository) UserServiceImpl {
    return UserServiceImpl{UserRepo: userRepo}
}