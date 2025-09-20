package ecies

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
    "log"

	"golang.org/x/crypto/hkdf"
)

type Envelope struct {
	Ciphertext   string `json:"ciphertext"`
	Nonce        string `json:"nonce"`
	EphemeralPub string `json:"ephemeralPub"`
}

var curve = elliptic.P256()

func Encrypt(pubKey *ecdsa.PublicKey , data []byte) ([]byte, error) {
    log.Printf("[ECIES] Encrypt: dataLen=%d", len(data))
    if pubKey != nil {
        log.Printf("[ECIES] Encrypt: recipient pub (P-256) X=%s Y=%s", pubKey.X.Text(16), pubKey.Y.Text(16))
    }
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
    log.Printf("[ECIES] Encrypt: eph pub X=%s Y=%s, ephPrivLen=%d", ex.Text(16), ey.Text(16), len(ephemeralPriv))
	// ECDH
	sx, sy := curve.ScalarMult(pubKey.X, pubKey.Y, ephemeralPriv)
	shared := elliptic.Marshal(curve,sx,sy)
    log.Printf("[ECIES] Encrypt: shared X=%s Y=%s, sharedBytesB64=%s", sx.Text(16), sy.Text(16), base64.StdEncoding.EncodeToString(shared))

	// HKDF
	h := hkdf.New(sha256.New, shared, nil, nil)
	aesKey := make([]byte, 32)
	_, err = io.ReadFull(h, aesKey)
	if err != nil {
		return nil, err
	}
    log.Printf("[ECIES] Encrypt: aesKeyB64=%s", base64.StdEncoding.EncodeToString(aesKey))

	// AES-GCM
	nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}
    log.Printf("[ECIES] Encrypt: nonceB64=%s", base64.StdEncoding.EncodeToString(nonce))
	blockAes, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blockAes)
	if err != nil {
		return nil, err
	}
	ct := gcm.Seal(nil, nonce, data, nil)
    log.Printf("[ECIES] Encrypt: ctLen=%d ctB64=%s", len(ct), base64.StdEncoding.EncodeToString(ct))

	ephemeralPubBytes, err := x509.MarshalPKIXPublicKey(ephemeralPub)
	if err != nil {
		return nil, err
	}
	env := Envelope{
		Ciphertext: base64.StdEncoding.EncodeToString(ct),
		Nonce: base64.StdEncoding.EncodeToString(nonce),
		EphemeralPub: base64.StdEncoding.EncodeToString(ephemeralPubBytes),
	}

	out, err := json.Marshal(env)
	if err == nil {
		log.Printf("[ECIES] Encrypt: envJSON=%s", string(out))
	}
	return out, err
}

func Decrypt(privKey *ecdsa.PrivateKey, env Envelope) ([]byte, error) {
    log.Printf("[ECIES] Decrypt: env.nonceB64=%s env.ciphertextLen=%d env.ephPubLen=%d",
        env.Nonce, len(env.Ciphertext), len(env.EphemeralPub))
	ep, err := base64.StdEncoding.DecodeString(env.EphemeralPub)
	if err != nil {
		return nil, err
	}
    log.Printf("[ECIES] Decrypt: ephPubDerLen=%d", len(ep))
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
    log.Printf("[ECIES] Decrypt: eph pub X=%s Y=%s", ePub.X.Text(16), ePub.Y.Text(16))
	
	// ECDH: shared = d_R * Q_e
	sx, sy := curve.ScalarMult(ePub.X, ePub.Y, privKey.D.Bytes())
	shared := elliptic.Marshal(curve,sx,sy)
    log.Printf("[ECIES] Decrypt: shared X=%s Y=%s, sharedBytesB64=%s", sx.Text(16), sy.Text(16), base64.StdEncoding.EncodeToString(shared))

	// HKDF
	h := hkdf.New(sha256.New, shared, nil, nil)
	aesKey := make([]byte, 32)
	_, err = io.ReadFull(h, aesKey)
	if err != nil {
		return nil, err
	}
    log.Printf("[ECIES] Decrypt: aesKeyB64=%s", base64.StdEncoding.EncodeToString(aesKey))

	// AES-GCM
	nonce, err := base64.StdEncoding.DecodeString(env.Nonce)
	if err != nil {
		return nil, err
	}
	ct, err := base64.StdEncoding.DecodeString(env.Ciphertext)
	if err != nil {
		return nil, err
	}
    log.Printf("[ECIES] Decrypt: nonceLen=%d ctLen=%d", len(nonce), len(ct))
	blockAes, err := aes.NewCipher(aesKey)
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
    log.Printf("[ECIES] Decrypt: ptLen=%d ptPreview=%q", len(pt), string(pt))
	return pt, nil
}