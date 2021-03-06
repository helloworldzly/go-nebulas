// Copyright (C) 2017 go-nebulas authors
//
// This file is part of the go-nebulas library.
//
// the go-nebulas library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-nebulas library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-nebulas library.  If not, see <http://www.gnu.org/licenses/>.
//

package account

import (
	"errors"

	"github.com/nebulasio/go-nebulas/core"
	"github.com/nebulasio/go-nebulas/crypto"
	"github.com/nebulasio/go-nebulas/crypto/cipher"
	"github.com/nebulasio/go-nebulas/crypto/keystore"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrTxAddressLocked from address locked.
	ErrTxAddressLocked = errors.New("transaction from address locked")
)

// Manager accounts manager ,handle account generate and storage
type Manager struct {
	ks *keystore.Keystore
}

// NewManager new a account manager
func NewManager() *Manager {
	m := new(Manager)
	m.ks = keystore.DefaultKS
	return m
}

// NewAccount returns a new address and keep it in keystore
func (m *Manager) NewAccount(passphrase []byte) (*core.Address, error) {

	priv, err := crypto.NewPrivateKey(keystore.SECP256K1, nil)
	if err != nil {
		return nil, err
	}
	return m.storeAddress(priv, passphrase)
}

func (m *Manager) storeAddress(priv keystore.PrivateKey, passphrase []byte) (*core.Address, error) {
	pub, err := priv.PublicKey().Encoded()
	if err != nil {
		return nil, err
	}
	addr, err := core.NewAddressFromPublicKey(pub)
	if err != nil {
		return nil, err
	}

	err = m.ks.SetKey(addr.ToHex(), priv, passphrase)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// Unlock unlock address with passphrase
func (m *Manager) Unlock(addr *core.Address, passphrase []byte) error {
	return m.ks.Unlock(addr.ToHex(), passphrase, keystore.DefaultUnlockDuration)
}

// Lock lock address
func (m *Manager) Lock(addr *core.Address) error {
	return m.ks.Lock(addr.ToHex())
}

// Accounts returns slice of address
func (m *Manager) Accounts() []*core.Address {
	aliases := m.ks.Aliases()
	addres := make([]*core.Address, len(aliases))
	for _, a := range aliases {
		addr, err := core.AddressParse(a)
		if err == nil {
			// currently keystore only storage address as alias
			addres = append(addres, addr)
		}
	}
	return addres
}

// Import import a key file to keystore, compatible ethereum keystore file
func (m *Manager) Import(keyjson, passphrase []byte) (*core.Address, error) {
	cipher := cipher.NewCipher(uint8(keystore.SCRYPT))
	data, err := cipher.DecryptKey(keyjson, passphrase)
	if err != nil {
		return nil, err
	}
	priv, err := crypto.NewPrivateKey(keystore.SECP256K1, data)
	if err != nil {
		return nil, err
	}
	return m.storeAddress(priv, passphrase)
}

// Export export address to key file
func (m *Manager) Export(addr *core.Address, passphrase []byte) ([]byte, error) {
	key, err := m.ks.GetKey(addr.ToHex(), passphrase)
	if err != nil {
		return nil, err
	}
	data, err := key.Encoded()
	if err != nil {
		return nil, err
	}
	cipher := cipher.NewCipher(uint8(keystore.SCRYPT))
	if err != nil {
		return nil, err
	}
	out, err := cipher.EncryptKey(addr.ToHex(), data, passphrase)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SignTransaction sign transaction with the specified algorithm
func (m *Manager) SignTransaction(addr *core.Address, tx *core.Transaction) error {
	// TODO(larry.wang): check the addr is the tx's from address
	key, err := m.ks.GetUnlocked(addr.ToHex())
	if err != nil {
		log.WithFields(log.Fields{
			"func": "SignTransaction",
			"err":  ErrTxAddressLocked,
			"tx":   tx,
		}).Error("transaction address locked")
		return err
	}

	signature, err := crypto.NewSignature(keystore.SECP256K1)
	if err != nil {
		return err
	}
	signature.InitSign(key.(keystore.PrivateKey))
	return tx.Sign(signature)
}
