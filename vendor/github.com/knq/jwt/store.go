package jwt

import (
	"crypto"
	"sync"
)

// Store is the common interface for a keystore.
type Store interface {
	// PublicKey returns the public key for a store.
	PublicKey() (crypto.PublicKey, bool)

	// PrivateKey returns the private key for a store.
	PrivateKey() (crypto.PrivateKey, bool)
}

// Keystore is a simple type providing a Store implementation.
type Keystore struct {
	// Key is the private key.
	Key interface{}

	// PublicKey is the public key.
	PubKey interface{}

	rw sync.RWMutex
}

// PublicKey returns the stored public key for the keystore, alternately
// generating the public key from the private key if the public key was not
// supplied and the private key was.
func (ks *Keystore) PublicKey() (crypto.PublicKey, bool) {
	ks.rw.RLock()
	key, pub := ks.Key, ks.PubKey
	ks.rw.RUnlock()
	if pub != nil {
		return pub, true
	}

	// generate the public key
	if key != nil {
		ks.rw.Lock()
		defer ks.rw.Unlock()

		if x, ok := key.(interface {
			Public() crypto.PublicKey
		}); ok {
			ks.PubKey = x.Public()
		}

		return ks.PubKey, ks.PubKey != nil
	}

	return nil, false
}

// PrivateKey returns the stored private key for the keystore.
func (ks *Keystore) PrivateKey() (crypto.PrivateKey, bool) {
	ks.rw.RLock()
	defer ks.rw.RUnlock()
	return ks.Key, ks.Key != nil
}
