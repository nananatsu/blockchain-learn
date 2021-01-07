package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

// Wallet 钱包
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

const version = byte(0x00)
const addressChecksumLen = 4

// GetAddress 获取钱包地址
func (w *Wallet) GetAddress() []byte {

	pubKeyHash := HashPubKey(w.PublicKey)

	version := append([]byte{version}, pubKeyHash...)
	checksum := checkSum(version)

	payload := append(version, checksum...)
	address := Base58Encode(payload)

	return address

}

// ValidateAddress 检查地址是否合法
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))

	actualCheckSum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetCheckSUm := checkSum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualCheckSum, targetCheckSUm) == 0
}

// HashPubKey 计算pubkey hash
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()

	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

func checkSum(payload []byte) []byte {
	sha1 := sha256.Sum256(payload)
	sha2 := sha256.Sum256(sha1[:])

	return sha2[:addressChecksumLen]
}

// NewWallet 新建钱包
func NewWallet() *Wallet {
	private, public := newKeyPair()

	wallet := Wallet{private, public}

	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {

	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		log.Panic(err)
	}
	publicKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, publicKey
}
