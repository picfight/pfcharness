// Copyright (c) 2018 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package memwallet

import (
	"github.com/picfight/pfcd/chaincfg"
	"github.com/picfight/pfcd/pfcec"
	"github.com/picfight/pfcd/wire"
	"github.com/picfight/pfcutil"
)

const chainUpdateSignal = "chainUpdateSignal"
const stopSignal = "stopSignal"

// chainUpdate encapsulates an update to the current main chain. This struct is
// used to sync up the InMemoryWallet each time a new block is connected to the main
// chain.
type chainUpdate struct {
	blockHeight  int32
	filteredTxns []*pfcutil.Tx
	isConnect    bool // True if connect, false if disconnect
}

// undoEntry is functionally the opposite of a chainUpdate. An undoEntry is
// created for each new block received, then stored in a log in order to
// properly handle block re-orgs.
type undoEntry struct {
	utxosDestroyed map[wire.OutPoint]*utxo
	utxosCreated   []wire.OutPoint
}

// utxo represents an unspent output spendable by the InMemoryWallet. The maturity
// height of the transaction is recorded in order to properly observe the
// maturity period of direct coinbase outputs.
type utxo struct {
	pkScript       []byte
	value          pfcutil.Amount
	maturityHeight int32
	keyIndex       uint32
	isLocked       bool
}

// isMature returns true if the target utxo is considered "mature" at the
// passed block height. Otherwise, false is returned.
func (u *utxo) isMature(height int32) bool {
	return height >= u.maturityHeight
}

// keyToAddr maps the passed private to corresponding p2pkh address.
func keyToAddr(key *pfcec.PrivateKey, net *chaincfg.Params) (pfcutil.Address, error) {
	serializedKey := key.PubKey().SerializeCompressed()
	pubKeyAddr, err := pfcutil.NewAddressPubKey(serializedKey, net)
	if err != nil {
		return nil, err
	}
	return pubKeyAddr.AddressPubKeyHash(), nil
}
