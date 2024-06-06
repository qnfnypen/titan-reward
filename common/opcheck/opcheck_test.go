package opcheck

import (
	"encoding/hex"
	"fmt"
	"testing"

	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestVerifyAddrSign(t *testing.T) {
	// key := "02AE577C704CA710D2083AEB86ABB0C97840A9E276A6D3E026485D921BB62C2A6B"
	key := "02880afd856ceaf0ea642fb34d8654652260dc40dc7bc38ac458b50ca0e59d495f"

	pr := "titan1jr4def3jn7a6x2kn7klt638w9xfuxuf8zjala7"

	// titan1jr4def3jn7a6x2kn7klt638w9xfuxuf8zjala7

	match, err := VerifyComosAddr(pr, key, "titan")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(match)
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
	t.Log(address)
	t.Log(pubKey.String())
	bs, _ := hex.DecodeString("034AFC08B13BF8CCB8681DF260CF8D12B05F15A5F7CAC30422C8E560A57E587D00")
	t.Log(fmt.Sprintf("%x", bs))
	t.Log(fmt.Sprintf("%x", pubKey.Bytes()))
	t.Logf("Public Key: %s\n", hex.EncodeToString(pubKey.Bytes()))
	t.Logf("Address: %s\n", address.String())
}
