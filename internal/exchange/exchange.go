package exchange

import (
	"context"
	"crypto/ecdh"
	cryptoRand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
)

type keyStore interface {
	Set(ctx context.Context, key string, value []byte) error

	Get(ctx context.Context, key string) ([]byte, error)
}

var store keyStore
var curve ecdh.Curve
var rand io.Reader

func init() {
	store = newInmemoryStore()
	curve = ecdh.P256()
	rand = cryptoRand.Reader
}

func Exchange(ctx context.Context, encodedPeerPubKey []byte) (string, []byte, error) {
	decodedPeerPubKey := make([]byte, hex.DecodedLen(len(encodedPeerPubKey)))
	if _, err := hex.Decode(decodedPeerPubKey, encodedPeerPubKey); err != nil {
		return "", nil, fmt.Errorf("decode public key: %w", err)
	}

	privKey, err := curve.GenerateKey(rand)
	if err != nil {
		return "", nil, fmt.Errorf("generate private key: %w", err)
	}

	peerPubKey, err := curve.NewPublicKey(decodedPeerPubKey)
	if err != nil {
		return "", nil, fmt.Errorf("generate public key: %w", err)
	}

	secret, err := privKey.ECDH(peerPubKey)
	if err != nil {
		return "", nil, fmt.Errorf("exchange: %w", err)
	}

	encodedSecret := make([]byte, hex.EncodedLen(len(secret)))
	hex.Encode(encodedSecret, secret)

	sessionID := uuid.NewString()
	if err = store.Set(ctx, sessionID, encodedSecret); err != nil {
		return "", nil, err
	}

	decodedPubKey := privKey.PublicKey().Bytes()
	encodedPubKey := make([]byte, hex.EncodedLen(len(decodedPeerPubKey)))
	hex.Encode(encodedPubKey, decodedPubKey)

	return sessionID, encodedPubKey, nil
}

var ErrSecretNotFound = errors.New("secret not found")

func GetSecret(ctx context.Context, sessionID string) ([]byte, error) {
	secret, err := store.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
