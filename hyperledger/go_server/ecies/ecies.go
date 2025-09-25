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
		WrappedKeys  map[string]string `json:"wrappedKeys"` // label -> b64(smallWrapJSON)
		Enc          string            `json:"enc"`         // "AES-256-GCM"
		KDF         string            `json:"kdf"`         // "HKDF-SHA256"
		Curve       string            `json:"curve"`       // "P-256"
	}

	// smallWrap：包 dataKey 的小包（再 b64 放進 WrappedKeys）
	type smallWrap struct {
		EphPub string `json:"ephPub"` // b64(未壓縮65B)
		N      string `json:"n"`      // b64(12B)
		C      string `json:"c"`      // b64(AES-GCM(kek, dataKey))
		KDF    string `json:"kdf"`    // "HKDF-SHA256"
		Curve  string `json:"curve"`  // "P-256"
	}

	var (
		curve = elliptic.P256()
		b64   = base64.StdEncoding
	)


// EncryptReportForClinic：A/B/C 三步到位（同時包診所與平台兩份）
func EncryptReportForClinic(clinicPub *ecdsa.PublicKey, platformPub *ecdsa.PublicKey, plaintext []byte, clinicLabel string) ([]byte, error) {
		if clinicPub == nil || clinicLabel == "" {
			return nil, errors.New("clinic pub/label required")
		}
		if platformPub == nil {
			return nil, errors.New("platform pub required")
		}

		log.Printf("[ECIES] EncryptReportForClinic: clinicLabel=%s, ptLen=%d", clinicLabel, len(plaintext))
		if clinicPub != nil {
			log.Printf("[ECIES] ClinicPub: X=%s Y=%s", clinicPub.X.Text(16), clinicPub.Y.Text(16))
		}
		if platformPub != nil {
			log.Printf("[ECIES] PlatformPub: X=%s Y=%s", platformPub.X.Text(16), platformPub.Y.Text(16))
		}

		// A) dataKey = 32B 隨機；nonce = 12B
		dataKey := make([]byte, 32)
		if _, err := io.ReadFull(rand.Reader, dataKey); err != nil { return nil, err }
		log.Printf("[ECIES] dataKeyB64=%s", b64.EncodeToString(dataKey))
		nonce := make([]byte, 12)
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil { return nil, err }
		log.Printf("[ECIES] dataNonceB64=%s", b64.EncodeToString(nonce))

		// B) AES-GCM(dataKey) 加密報告
		block, err := aes.NewCipher(dataKey)
		if err != nil { return nil, err }
		gcm, err := cipher.NewGCM(block)
		if err != nil { return nil, err }
		ct := gcm.Seal(nil, nonce, plaintext, nil)
		log.Printf("[ECIES] dataCTLen=%d dataCTB64=%s", len(ct), b64.EncodeToString(ct))

		// C) 為診所與平台包 dataKey（ECIES：ephemeral×pub → HKDF → kek → GCM 包 dataKey）
		wrapClinicB64, err := wrapDataKeyFor(clinicPub, dataKey)
		if err != nil { return nil, err }
		log.Printf("[ECIES] wrap(for=%s) len=%d", clinicLabel, len(wrapClinicB64))

		wrapPlatformB64, err := wrapDataKeyFor(platformPub, dataKey)
		if err != nil { return nil, err }
		log.Printf("[ECIES] wrap(for=platform) len=%d", len(wrapPlatformB64))

		env := Envelope{
			Ciphertext:  b64.EncodeToString(ct),
			Nonce:       b64.EncodeToString(nonce),
			WrappedKeys: map[string]string{clinicLabel: wrapClinicB64, "platform": wrapPlatformB64},
			Enc:         "AES-256-GCM",
			KDF:         "HKDF-SHA256",
			Curve:       "P-256",
		}
		out, err := json.Marshal(env)
		if err == nil {
			log.Printf("[ECIES] EnvelopeJSON=%s", string(out))
		}
		return out, err
	}

	// wrapDataKeyFor：針對某收件者公鑰，產生該收件者專屬的小包
	func wrapDataKeyFor(recipientPub *ecdsa.PublicKey, dataKey []byte) (string, error) {
		if recipientPub == nil || recipientPub.Curve != curve || !curve.IsOnCurve(recipientPub.X, recipientPub.Y) {
			return "", errors.New("bad recipient pub")
		}
		log.Printf("[ECIES] wrapDataKeyFor: recPub X=%s Y=%s, dataKeyLen=%d", recipientPub.X.Text(16), recipientPub.Y.Text(16), len(dataKey))
		// ephemeral for this recipient
		ephPriv, ex, ey, err := elliptic.GenerateKey(curve, rand.Reader)
		if err != nil { return "", err }
		ephPub := &ecdsa.PublicKey{Curve: curve, X: ex, Y: ey}
		log.Printf("[ECIES] wrapDataKeyFor: ephPub X=%s Y=%s ephPrivLen=%d", ex.Text(16), ey.Text(16), len(ephPriv))

		// shared → kek
		sx, sy := curve.ScalarMult(recipientPub.X, recipientPub.Y, ephPriv)
		shared := elliptic.Marshal(curve, sx, sy)
		log.Printf("[ECIES] wrapDataKeyFor: shared X=%s Y=%s sharedB64=%s", sx.Text(16), sy.Text(16), b64.EncodeToString(shared))
		h := hkdf.New(sha256.New, shared, nil, []byte("ECIES-P256/WRAP-KEK/v1"))
		kek := make([]byte, 32)
		if _, err := io.ReadFull(h, kek); err != nil { return "", err }
		log.Printf("[ECIES] wrapDataKeyFor: kekB64=%s", b64.EncodeToString(kek))

		n := make([]byte, 12)
		if _, err := io.ReadFull(rand.Reader, n); err != nil { return "", err }
		log.Printf("[ECIES] wrapDataKeyFor: nB64=%s", b64.EncodeToString(n))
		block, err := aes.NewCipher(kek)
		if err != nil { return "", err }
		gcm, err := cipher.NewGCM(block)
		if err != nil { return "", err }
		wct := gcm.Seal(nil, n, dataKey, nil)
		log.Printf("[ECIES] wrapDataKeyFor: wctLen=%d wctB64=%s", len(wct), b64.EncodeToString(wct))

		epubDER, err := x509.MarshalPKIXPublicKey(ephPub)
		if err != nil { return "", err }
		sw := smallWrap{
			EphPub: b64.EncodeToString(epubDER),
			N:      b64.EncodeToString(n),
			C:      b64.EncodeToString(wct),
			KDF:    "HKDF-SHA256",
			Curve:  "P-256",
		}
		raw, err := json.Marshal(sw)
		if err != nil { return "", err }
		log.Printf("[ECIES] wrapDataKeyFor: smallWrapJSON=%s", string(raw))
		return b64.EncodeToString(raw), nil
	}

	// AddRecipient：用「現有持有人」取出 dataKey → 為新對象包一份 → 更新 Envelope
// AddRecipient：由平台呼叫。平台使用自身私鑰解開 label="platform" 的小包取得 dataKey，再為新對象包一份小包。
func AddRecipient(envJSON []byte, newLabel string, newRecipientPub *ecdsa.PublicKey,
		unwrapSelf func(env Envelope) ([]byte, error)) ([]byte, error) {

		if newLabel == "" || newRecipientPub == nil {
			return nil, errors.New("label/pub required")
		}
		var env Envelope
		if err := json.Unmarshal(envJSON, &env); err != nil { return nil, err }
		log.Printf("[ECIES] AddRecipient: newLabel=%s currentKeys=%d", newLabel, len(env.WrappedKeys))

		// 僅允許平台以自身私鑰解出 dataKey（label 必須為 "platform"）
		log.Printf("[ECIES] AddRecipient: unwrapSelf start (must be platform)")
		if _, ok := env.WrappedKeys["platform"]; !ok {
			return nil, errors.New("no platform wrapped key")
		}
		dataKey, err := unwrapSelf(env)
		if err != nil { return nil, err }
		log.Printf("[ECIES] AddRecipient: unwrapSelf done, dataKeyLen=%d", len(dataKey))

		// 包給新收件者（病患）
		wrapB64, err := wrapDataKeyFor(newRecipientPub, dataKey)
		if err != nil { return nil, err }
		if env.WrappedKeys == nil { env.WrappedKeys = make(map[string]string) }
		env.WrappedKeys[newLabel] = wrapB64
		out, err := json.Marshal(env)
		if err == nil {
			log.Printf("[ECIES] AddRecipient: EnvelopeJSON=%s", string(out))
		}
		return out, err
	}

	func DecryptReport(ownPriv *ecdsa.PrivateKey, env Envelope, ownLabel string) ([]byte, error) {
		log.Printf("[ECIES] DecryptReport: ownLabel=%s hasKeys=%d", ownLabel, len(env.WrappedKeys))
		wrapB64, ok := env.WrappedKeys[ownLabel]
		if !ok { return nil, errors.New("no wrappedKey for label") }

		// 拆小包得 dataKey
		raw, err := b64.DecodeString(wrapB64)
		if err != nil { return nil, err }
		var sw smallWrap
		if err := json.Unmarshal(raw, &sw); err != nil { return nil, err }
		log.Printf("[ECIES] DecryptReport: smallWrap=%s", string(raw))
		epb, err := b64.DecodeString(sw.EphPub)
		if err != nil { return nil, err }
		apub, err := x509.ParsePKIXPublicKey(epb); if err != nil { return nil, err }
		epub, ok := apub.(*ecdsa.PublicKey); if !ok { return nil, errors.New("bad eph pub type") }
		log.Printf("[ECIES] DecryptReport: ephPub X=%s Y=%s", epub.X.Text(16), epub.Y.Text(16))

		sx, sy := curve.ScalarMult(epub.X, epub.Y, ownPriv.D.Bytes())
		shared := elliptic.Marshal(curve, sx, sy)
		h := hkdf.New(sha256.New, shared, nil, []byte("ECIES-P256/WRAP-KEK/v1"))
		kek := make([]byte, 32)
		if _, err := io.ReadFull(h, kek); err != nil { return nil, err }
		log.Printf("[ECIES] DecryptReport: kekB64=%s", b64.EncodeToString(kek))

		n, err := b64.DecodeString(sw.N)
		if err != nil { return nil, err }
		c, err := b64.DecodeString(sw.C)
		if err != nil { return nil, err }
		log.Printf("[ECIES] DecryptReport: wrap nLen=%d cLen=%d", len(n), len(c))
		block, err := aes.NewCipher(kek)
		if err != nil { return nil, err }
		gcm, err := cipher.NewGCM(block)
		if err != nil { return nil, err }
		dataKey, err := gcm.Open(nil, n, c, nil)
		if err != nil { return nil, err }
		log.Printf("[ECIES] DecryptReport: dataKeyLen=%d dataKeyB64=%s", len(dataKey), b64.EncodeToString(dataKey))

		// 用 dataKey 解主密文
		nonce, err := b64.DecodeString(env.Nonce)
		if err != nil { return nil, err }
		ct, err := b64.DecodeString(env.Ciphertext)
		if err != nil { return nil, err }
		log.Printf("[ECIES] DecryptReport: main nonceLen=%d ctLen=%d", len(nonce), len(ct))
		block2, err := aes.NewCipher(dataKey)
		if err != nil { return nil, err }
		gcm2, err := cipher.NewGCM(block2)
		if err != nil { return nil, err }
		pt, err := gcm2.Open(nil, nonce, ct, nil)
		if err != nil { return nil, err }
		log.Printf("[ECIES] DecryptReport: ptLen=%d ptPreview=%q", len(pt), string(pt))
		return pt, nil
	}
