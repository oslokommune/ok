package jit

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	"github.com/zalando/go-keyring"
)

const (
	keyringService = "ok-jit"
	keyringUser    = "encryption-key"
	cacheFileName  = "jit-msal-cache.enc"
)

// tokenCache implements the MSAL cache.ExportReplace interface.
// Data is encrypted on disk with AES-GCM; the key is stored in the OS keychain.
type tokenCache struct {
	path string
}

func newTokenCache() (*tokenCache, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(home, ".config", "ok")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	return &tokenCache{path: filepath.Join(dir, cacheFileName)}, nil
}

func (t *tokenCache) Replace(ctx context.Context, unmarshaler cache.Unmarshaler, hints cache.ReplaceHints) error {
	encrypted, err := os.ReadFile(t.path)
	if err != nil {
		if os.IsNotExist(err) {
			// No cache file yet, nothing to load
			return nil
		}
		return fmt.Errorf("reading cache file: %w", err)
	}

	key, err := getOrCreateKey()
	if err != nil {
		return fmt.Errorf("loading encryption key: %w", err)
	}

	data, err := decrypt(key, encrypted)
	if err != nil {
		return fmt.Errorf("decrypting cache (file may be tampered or key rotated): %w", err)
	}

	return unmarshaler.Unmarshal(data)
}

func (t *tokenCache) Export(ctx context.Context, marshaler cache.Marshaler, hints cache.ExportHints) error {
	data, err := marshaler.Marshal()
	if err != nil {
		return err
	}

	key, err := getOrCreateKey()
	if err != nil {
		return err
	}

	encrypted, err := encrypt(key, data)
	if err != nil {
		return err
	}

	return os.WriteFile(t.path, encrypted, 0600)
}

// ClearCache removes all cached tokens from disk and keychain.
func ClearCache() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	os.Remove(filepath.Join(home, ".config", "ok", cacheFileName))
	keyring.Delete(keyringService, keyringUser)

	return nil
}

// RefreshTokenInfo contains info about cached refresh tokens.
type RefreshTokenInfo struct {
	Present bool
	Secret  string
}

// GetRefreshTokenInfo checks the encrypted cache for refresh token details.
func GetRefreshTokenInfo() RefreshTokenInfo {
	home, err := os.UserHomeDir()
	if err != nil {
		return RefreshTokenInfo{}
	}

	encrypted, err := os.ReadFile(filepath.Join(home, ".config", "ok", cacheFileName))
	if err != nil {
		return RefreshTokenInfo{}
	}

	key, err := getOrCreateKey()
	if err != nil {
		return RefreshTokenInfo{}
	}

	data, err := decrypt(key, encrypted)
	if err != nil {
		return RefreshTokenInfo{}
	}

	var cacheData struct {
		RefreshToken map[string]struct {
			Secret string `json:"secret"`
		} `json:"RefreshToken"`
	}
	if json.Unmarshal(data, &cacheData) != nil {
		return RefreshTokenInfo{}
	}

	for _, rt := range cacheData.RefreshToken {
		return RefreshTokenInfo{Present: true, Secret: rt.Secret}
	}

	return RefreshTokenInfo{}
}

func getOrCreateKey() ([]byte, error) {
	keyHex, err := keyring.Get(keyringService, keyringUser)
	if err == nil {
		return hex.DecodeString(keyHex)
	}

	// Generate a new 256-bit key
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	if err := keyring.Set(keyringService, keyringUser, hex.EncodeToString(key)); err != nil {
		return nil, err
	}

	return key, nil
}

func encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	return gcm.Open(nil, ciphertext[:nonceSize], ciphertext[nonceSize:], nil)
}
