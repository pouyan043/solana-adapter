package solanaadapter

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"

	"github.com/mr-tron/base58"
	log "github.com/sirupsen/logrus"
	"github.com/stellar/go/exp/crypto/derivation"
)

type SolanaAdapter struct {
	logger *log.Logger
}

func NewSolanaAdapter(logger *log.Logger) *SolanaAdapter {
	return &SolanaAdapter{logger: logger}
}

func (adapter *SolanaAdapter) CanDo(coinType uint) bool {
	return coinType == 501
}

func (adapter *SolanaAdapter) deriveKeysForPath(seed []byte, derivationPath string) (ed25519.PrivateKey, ed25519.PublicKey, error) {
	derivedKey, err := derivation.DeriveForPath(derivationPath, seed)
	if err != nil {
		return nil, nil, err
	}
	privKey := ed25519.NewKeyFromSeed(derivedKey.Key[:])
	pubKey, ok := privKey.Public().(ed25519.PublicKey)
	if !ok {
		return nil, nil, errors.New("failed to get public key")
	}
	return privKey, pubKey, nil
}

func (adapter *SolanaAdapter) DerivePrivateKey(seed []byte, derivationPath string, isDev bool) (string, error) {
	privKey, _, err := adapter.deriveKeysForPath(seed, derivationPath)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(privKey), nil
}

func (adapter *SolanaAdapter) DerivePublicKey(seed []byte, derivationPath string, isDev bool) (string, error) {
	_, pubKey, err := adapter.deriveKeysForPath(seed, derivationPath)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(pubKey), nil
}

func (adapter *SolanaAdapter) DeriveAddress(seed []byte, derivationPath string, isDev bool) (string, error) {
	_, pubKey, err := adapter.deriveKeysForPath(seed, derivationPath)
	if err != nil {
		return "", err
	}
	return base58.Encode(pubKey), nil
}

func (adapter *SolanaAdapter) GetBlockchainNetwork(isDev bool) string {
	if isDev {
		return "dev"
	}
	return "mainnet"
}

func (adapter *SolanaAdapter) CreateSignedTransaction(seed []byte, derivationPath string, payload string) (string, error) {
	privKey, _, err := adapter.deriveKeysForPath(seed, derivationPath)
	if err != nil {
		return "", err
	}
	message, err := hex.DecodeString(payload)
	if err != nil {
		return "", err
	}
	signature := ed25519.Sign(privKey, message)
	return hex.EncodeToString(signature), nil
}

func (adapter *SolanaAdapter) CreateSignature(seed []byte, derivationPath string, payload string) (string, error) {
	return adapter.CreateSignedTransaction(seed, derivationPath, payload)
}
