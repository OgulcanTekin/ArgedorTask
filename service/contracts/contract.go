package contract

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/gin-gonic/gin"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	contract "example.com/smart-contract/contracts"
)

type DeploymentRequest struct {
	walletAddress string
}

func GetAuth(client *ethclient.Client, privateAddress string) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(privateAddress)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	blockchainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, blockchainId)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = big.NewInt(1000000)

	return auth, nil
}

func DeployContract(c *gin.Context) {
	// TODO: db sorgusu at ordan gelen private key 63deki privatewID ile eşle

	var req DeploymentRequest
	c.BindJSON(&req)

	privateWalletID := req.walletAddress

	if len(privateWalletID) <= 0 {
		c.JSON(401, gin.H{
			"error": "Wallet ID is empty",
		})
		return
	}
	client, err := ethclient.Dial("https://api.avax.network/ext/bc/C/rpc")
	if err != nil {
		c.JSON(401, gin.H{
			"error": err,
		})
		return
	}
	auth, err := GetAuth(client, privateWalletID)
	if err != nil {
		c.JSON(401, gin.H{
			"error": err,
		})
		return
	}
	depContAdd, tx, instance, err := contract.DeployErc20(auth, client, "Task-Test", "1")
	if err != nil {
		c.JSON(401, gin.H{
			"error": err,
		})
		return
	}
	fmt.Println(depContAdd.Hex(), tx, instance)
	c.JSON(401, gin.H{
		"message": "Contract Successfully Deployed",
	})
}

func transfer() {
	// transfer function yazilmasi
}

func mint() {
	// Mint istenilen sekilde eklenmesi
}
func getBalance() {
	// Balance function yazilmasi
}
