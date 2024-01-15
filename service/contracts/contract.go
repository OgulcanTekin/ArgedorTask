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

	// Veritabanı bağlantısını kurmak için kullanılan yapı
	// "github.com/go-sql-driver/postgresql" // Veritabanı sürücüsü
	// "github.com/jmoiron/sqlx"

	contract "example.com/smart-contract/contracts"
	// "example.com/smart-contract/service/wallet"
)

type DeploymentRequest struct {
	walletAddress string
}
type TransferRequest struct {
	SenderID   int     `json:"sender_id"`
	ReceiverID int     `json:"receiver_id"`
	Amount     float64 `json:"amount"`
}

// Veritabanı bağlantısını kurmak için kullanılan yapı
// type Database struct {
// 	Conn *sqlx.DB
// }

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
	// TODO: db sorgusu ile tablodan gelen private key 63deki privateWalletID ile eşle

	var req DeploymentRequest
	c.BindJSON(&req)

	// private_key, err := getPublicKeyByPrivateWalletID(privateWalletID)
	// if err != nil {
	// 	c.JSON(401, gin.H{
	// 		"error": "Error retrieving wallet ID from the database",
	// 	})
	// 	return
	// }
	//id.walletID
	privateWalletID := req.walletAddress

	if len(privateWalletID) <= 0 {
		c.JSON(401, gin.H{
			"error": "Wallet ID is empty",
		})
		return
	}

	// private_key, err := d.getPrivateKeyByPrivateWalletID(privateWalletID)
	// if err != nil {
	// 	c.JSON(401, gin.H{
	// 		"error": "Error retrieving wPrivate Key from the database",
	// 	})
	// 	return
	// }

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

// func getPrivateKeyByPrivateWalletID(private_key string) (string, error) {
// 	// TODO: Veritabanına sorgu gönder ve privateWalletID'ye karşılık gelen private_key'yi al
// 	var private_key string
// 	err := db.QueryRow("SELECT id FROM wallets WHERE private_key = ?", id).Scan(&walletID)
// 	return private_key, err

//		return 0, nil
//	}

// func (d *Database) Transfer(c *gin.Context) {
// 	var req TransferRequest
// 	c.BindJSON(&req)

// 	senderID := req.SenderID
// 	receiverID := req.ReceiverID
// 	amount := req.Amount

// 	if senderID == receiverID {
// 		c.JSON(400, gin.H{
// 			"error": "Sender and receiver IDs should be different",
// 		})
// 		return
// 	}

// 	// Gönderen cüzdanın bakiyesini kontrol et
// 	senderWallet, err := d.getWalletByID(senderID)
// 	if err != nil {
// 		c.JSON(500, gin.H{
// 			"error": "Error retrieving sender wallet information",
// 		})
// 		return
// 	}

// 	if senderWallet.Balance < amount {
// 		c.JSON(400, gin.H{
// 			"error": "Insufficient balance",
// 		})
// 		return
// 	}

// 	// Alıcı cüzdanın bakiyesini güncelle
// 	receiverWallet, err := d.getWalletByID(receiverID)
// 	if err != nil {
// 		c.JSON(500, gin.H{
// 			"error": "Error retrieving receiver wallet information",
// 		})
// 		return
// 	}

// 	// İşlemleri başarılı bir şekilde gerçekleştir
// 	err = d.executeTransaction(senderID, receiverID, amount)
// 	if err != nil {
// 		c.JSON(500, gin.H{
// 			"error": "Error executing transaction",
// 		})
// 		return
// 	}

// 	c.JSON(200, gin.H{
// 		"message": "Transaction successfully completed",
// 	})

// }

func mint() {
	// Mint istenilen sekilde eklenmesi
}
func getBalance() {
	// Balance function yazilmasi
}
