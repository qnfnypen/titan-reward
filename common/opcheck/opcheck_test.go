package opcheck

import (
	"encoding/hex"
	"fmt"
	"testing"

	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestVerifyAddrSign(t *testing.T) {
	nonce := "Hello, Cosmos!"
	sign := "868177c2f4a206fea2cfa90ea94d2cc49a1b0cdf27d7a29a34ebbcc498568b8d585314f4b660bf954fb6d7a76a021c5cf9fa781d56895060ab8e01969394a942"
	// key := "02C90E04DCFCA1A8B4A9CC40223F8015D9A4665409F33B84ACA8E48178044B47A9"
	key := "034AFC08B13BF8CCB8681DF260CF8D12B05F15A5F7CAC30422C8E560A57E587D00"

	bs, err := hex.DecodeString(key)
	if err != nil {
		t.Fatal(err)
	}

	pubkey := &secp256k1.PubKey{Key: bs}
	b, err := VerifySignature(nonce, sign, pubkey)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(b)
}

func TestComos(t *testing.T) {
	// 生成新的私钥
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()

	// 要签名的消息
	message := "Hello, Cosmos!"

	// 对消息进行签名
	signature, err := SignMessage(message, privKey)
	if err != nil {
		t.Fatalf("Failed to sign message: %v", err)
	}
	t.Logf("Message: %s\n", message)
	t.Logf("Signature: %s\n", signature)

	// 验证签名
	valid, err := VerifySignature(message, signature, pubKey)
	if err != nil {
		t.Fatalf("Failed to verify signature: %v", err)
	}

	if valid {
		t.Logf("Signature is valid")
	} else {
		t.Logf("Signature is invalid")
	}

	// 打印生成的公钥和地址
	address := sdk.AccAddress(pubKey.Address())
	t.Log(pubKey.String())
	bs, _ := hex.DecodeString("034AFC08B13BF8CCB8681DF260CF8D12B05F15A5F7CAC30422C8E560A57E587D00")
	t.Log(fmt.Sprintf("%x", bs))
	t.Log(fmt.Sprintf("%x", pubKey.Bytes()))
	t.Logf("Public Key: %s\n", hex.EncodeToString(pubKey.Bytes()))
	t.Logf("Address: %s\n", address.String())
}
