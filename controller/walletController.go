package controller

import (
	"example.com/smart-contract/db"
	contract "example.com/smart-contract/service/contracts"
	wallet "example.com/smart-contract/service/wallet"
	"github.com/gin-gonic/gin"
)

func StartService() {
	service := gin.Default()

	db.InitializeDatabase()

	service.POST("/createWallet", wallet.CreateWallet)
	service.POST("/deployContract", contract.DeployContract) // Buraya contract.go üzerinde yazılan methodlar gelecek deploy transfer vs

	service.Run()
}
