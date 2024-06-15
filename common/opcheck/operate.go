package opcheck

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type signDocFee struct {
	Amount []sdk.Coin `json:"amount"`
	Gas    string     `json:"gas"`
}

type signDocMsgValue struct {
	Data   string `json:"data"`
	Signer string `json:"signer"`
}

type signDocMsg struct {
	Type  string          `json:"type"`
	Value signDocMsgValue `json:"value"`
}
type signDoc struct {
	AccountNumber string       `json:"account_number"`
	ChainID       string       `json:"chain_id"`
	Fee           signDocFee   `json:"fee"`
	Memo          string       `json:"memo"`
	Msgs          []signDocMsg `json:"msgs"`
	Sequence      string       `json:"sequence"`
}

// ComposeArbitraryMsg Creates SignDoc with JSON encoded bytes as per adr036
// Compatible with AMINO as it is supported by keplr wallet
func ComposeArbitraryMsg(signer string, data string) ([]byte, error) {
	dataBase64 := base64.StdEncoding.EncodeToString([]byte(data))

	newSignDocMsgValue := signDocMsgValue{
		Data:   dataBase64,
		Signer: signer,
	}

	newSignDocMsg := signDocMsg{
		Value: newSignDocMsgValue,
		Type:  "sign/MsgSignData",
	}

	newSignDoc := signDoc{
		Msgs: []signDocMsg{
			newSignDocMsg,
		},
		AccountNumber: "0",
		Sequence:      "0",
		Fee: signDocFee{
			Gas:    "0",
			Amount: sdk.NewCoins(),
		},
	}

	jsonBytes, err := json.Marshal(newSignDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to Sign Doc to JSON: %w", err)
	}
	return jsonBytes, nil
}

// SignMessage signs a message using the provided private key
func SignMessage(message string, privKey cryptotypes.PrivKey) (string, error) {
	hash := sha256.Sum256([]byte(message))
	signature, err := privKey.Sign(hash[:])
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(signature), nil
}

// VerifySignature verifies a signed message using the provided public key
func VerifySignature(message, signatureHex string, pubKey cryptotypes.PubKey) (bool, error) {
	hash := sha256.Sum256([]byte(message))
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, err
	}
	return pubKey.VerifySignature(hash[:], signature), nil
}
