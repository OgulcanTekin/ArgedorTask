package wallet

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"example.com/smart-contract/db"
	walletModal "example.com/smart-contract/repository/modals"
)

func CreateWallet(c *gin.Context) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		c.JSON(401, gin.H{
			"error": err,
		})
		return
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println("Here is your private key. Do not share it with anyone!" + hexutil.Encode(privateKeyBytes)[2:])

	lastPrivateKey := hexutil.Encode(privateKeyBytes)[2:]

	db := db.Connection()

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		c.JSON(401, gin.H{
			"error": err,
		})
		return
	}

	userWallet := &walletModal.Wallet{
		PrivateKey: lastPrivateKey,
		PublicKey:  publicKeyECDSA.X.String(),
	}

	_, err = db.Model(userWallet).Insert()
	if err != nil {
		c.JSON(401, gin.H{
			"error": err,
		})
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("Public Key...:", hexutil.Encode(publicKeyBytes))

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("Address...:", address)

	c.JSON(200, gin.H{
		"Private key...:": lastPrivateKey,
	})
}
