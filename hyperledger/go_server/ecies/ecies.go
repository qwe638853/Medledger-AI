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
	"errors"
	"io"
)

type Envelope struct {
	Ciphertext   string `json:"ciphertext"`
	Nonce        string `json:"nonce"`
	EphemeralPub string `json:"ephemeralPub"`
}

var curve = elliptic.P256()

func Encrypt(pubKey *ecdsa.PublicKey , data []byte) ([]byte, error) {
	// 檢查公鑰是否有效
	if pubKey == nil || pubKey.Curve != curve || !curve.IsOnCurve(pubKey.X, pubKey.Y) {
		return nil, errors.New("invalid recipient public key (need ECDSA P-256)")
	}
	// 生成 ephemeral key
	ephemeralPriv, ex, ey, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
    // 生成 ephemeral public key
	ephemeralPub := &ecdsa.PublicKey{
		Curve: curve,
		X: ex,
		Y: ey,
	}
	// ECDH
	sx, sy := curve.ScalarMult(pubKey.X, pubKey.Y, ephemeralPriv)
	shared := elliptic.Marshal(curve,sx,sy)

	// HKDF
	h := hkdf.New(sha256.New, shared, nil, nil)
	aseKey := make([]byte, 32)
	err = io.ReadFull(h, aseKey)
	if err != nil {
		return nil, err
	}

	// AES-GCM
	nonce := make([]byte, 12)
	err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}
	blockAes, err := aes.NewCipher(aseKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blockAes)
	if err != nil {
		return nil, err
	}
	ct := gcm.Seal(nil, nonce, data, nil)

	ephemeralPubBytes, err := x509.MarshalPKIXPublicKey(ephemeralPub)
	if err != nil {
		return nil, err
	}
	env := Envelope{
		Ciphertext: base64.StdEncoding.EncodeToString(ct),
		Nonce: base64.StdEncoding.EncodeToString(nonce),
		EphemeralPub: base64.StdEncoding.EncodeToString(ephemeralPubBytes),
	}
	
	return json.Marshal(env)
}

func Decrypt(privKey *ecdsa.PrivateKey, env Envelope) ([]byte, error) {
	ep, err = base64.StdEncoding.DecodeString(env.EphemeralPub)
	if err != nil {
		return nil, err
	}
	anyPub, err := x509.ParsePKIXPublicKey(ep)
	if err != nil {
		return nil, err
	}
	ePub, ok := anyPub.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid ephemeral public key (need ECDSA P-256)")
	}
	if ePub.Curve != curve || !curve.IsOnCurve(ePub.X, ePub.Y) {
		return nil, errors.New("invalid ephemeral public key (need ECDSA P-256)")
	}
	
	// ECDH: shared = d_R * Q_e
	sx, sy := curve.ScalarMult(ePub.X, ePub.Y, privKey.D.Bytes())
	shared := elliptic.Marshal(curve,sx,sy)

	// HKDF
	h := hkdf.New(sha256.New, shared, nil, nil)
	aseKey := make([]byte, 32)
	err = io.ReadFull(h, aseKey)
	if err != nil {
		return nil, err
	}

	// AES-GCM
	nonce, err := base64.StdEncoding.DecodeString(env.Nonce)
	if err != nil {
		return nil, err
	}
	ct, err := base64.StdEncoding.DecodeString(env.Ciphertext)
	if err != nil {
		return nil, err
	}
	blockAes, err := aes.NewCipher(aseKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blockAes)
	if err != nil {
		return nil, err
	}
	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, err
	}
	return pt, nil
}