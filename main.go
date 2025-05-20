package main

import (
	"encoding/hex"
	"fmt"

	"solana-project/solanaadapter"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	adapter := solanaadapter.NewSolanaAdapter(logger)
	seed, _ := hex.DecodeString("4e9b2b5f3d8a1c9f7e6b2a4d8c7e3f5a1b9c2d4e6f7a8b9c0d1e2f3a4b5c6d7")
	derivationPath := "m/44'/501'/0'/0'"
	isDev := false
	coinType := uint(501)

	if !adapter.CanDo(coinType) {
		logger.Error("Coin type not supported")
		return
	}

	privKey, err := adapter.DerivePrivateKey(seed, derivationPath, isDev)
	if err != nil {
		logger.Errorf("Failed to derive private key: %v", err)
		return
	}
	fmt.Printf("Private Key (hex): %s\n", privKey)

	pubKey, err := adapter.DerivePublicKey(seed, derivationPath, isDev)
	if err != nil {
		logger.Errorf("Failed to derive public key: %v", err)
		return
	}
	fmt.Printf("Public Key (hex): %s\n", pubKey)

	address, err := adapter.DeriveAddress(seed, derivationPath, isDev)
	if err != nil {
		logger.Errorf("Failed to derive address: %v", err)
		return
	}
	fmt.Printf("Address (base58): %s\n", address)

	network := adapter.GetBlockchainNetwork(isDev)
	fmt.Printf("Blockchain Network: %s\n", network)

	payload := "48656c6c6f20536f6c616e61"
	signature, err := adapter.CreateSignedTransaction(seed, derivationPath, payload)
	if err != nil {
		logger.Errorf("Failed to create signed transaction: %v", err)
		return
	}
	fmt.Printf("Signature (hex): %s\n", signature)

	signature2, err := adapter.CreateSignature(seed, derivationPath, payload)
	if err != nil {
		logger.Errorf("Failed to create signature: %v", err)
		return
	}
	fmt.Printf("Signature from CreateSignature (hex): %s\n", signature2)
}
