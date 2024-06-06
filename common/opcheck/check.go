package opcheck

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// VerifyAddrSign 判断以太坊地址的签名是否正确
func VerifyAddrSign(nonce, sign string) (string, error) {
	// Hash the unsigned message using EIP-191
	hashedMessage := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(nonce)) + nonce)
	hash := crypto.Keccak256Hash(hashedMessage)
	// Get the bytes of the signed message
	decodedMessage := hexutil.MustDecode(sign)
	// Handles cases where EIP-115 is not implemented (most wallets don't implement it)
	if decodedMessage[64] == 27 || decodedMessage[64] == 28 {
		decodedMessage[64] -= 27
	}
	// Recover a public key from the signed message
	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), decodedMessage)
	if err != nil {
		return "", err
	}

	return crypto.PubkeyToAddress(*sigPublicKeyECDSA).String(), nil
}

// VerifyComosSign 验证 comos 链的签名
func VerifyComosSign(nonce, sign, pk string) (bool, error) {
	bs, err := hex.DecodeString(pk)
	if err != nil {
		return false, fmt.Errorf("decode public key from pk error:%w", err)
	}

	pubkey := &secp256k1.PubKey{Key: bs}
	return VerifySignature(nonce, sign, pubkey)
}

// VerifyComosAddr 验证 comos 链的地址是否正确
func VerifyComosAddr(paddr, pk, pr string) (bool, error) {
	bs, err := hex.DecodeString(pk)
	if err != nil {
		return false, fmt.Errorf("decode public key from pk error:%w", err)
	}

	pubkey := &secp256k1.PubKey{Key: bs}
	addr, err := sdk.Bech32ifyAddressBytes(pr, pubkey.Address())
	if err != nil {
		return false, fmt.Errorf("get address error:%w", err)
	}

	return strings.EqualFold(paddr, addr), nil
}
