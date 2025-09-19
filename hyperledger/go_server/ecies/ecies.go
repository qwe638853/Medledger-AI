package ecies

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"math/big"
)

type Envelope struct {
	Ciphertext   string `json:"ciphertext"`
	Nonce        string `json:"nonce"`
	EphemeralPub string `json:"ephemeralPub"`
}

var curve = elliptic.P256()

func Encrypt(pubKey *ecdsa.PublicKey , data []byte) ([]byte, error) {
	// 生成 ephemeral key
	ephemeralPriv, ex, ey, err := elliptic.GenerateKey(elliptic.P256(), rand.Reader)
    // 生成 ephemeral public key
	ephemeralPub := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X: ex,
		Y: ey,
	}
	// ECDH
	sx, sy := curve.ScalarMult(pubKey.X, pubKey.Y, ephemeralPriv)
	shared := elliptic.Marshal(curve,sx,sy)

	h := hkdf.New(sha256.New, shared, nil, nil)
	asekey := make([]byte, 32)
	io.Readfull(h, asekey)

	nonce := make([]byte, 12)
	io.Readfull(rand.Reader, nonce)

	
	
}
