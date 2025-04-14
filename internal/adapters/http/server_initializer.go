package http

import (
    "github.com/AmiraliFarazmand/PTC_Task/internal/app"
)

func InitializeHTTPServer(purchaseService app.PurchaseServiceImpl, userService app.UserServiceImpl) *GinServer {
    return NewGinServer(purchaseService, userService)
}