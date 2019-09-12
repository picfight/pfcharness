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
	blockHeight  int64
	filteredTxns []*pfcutil.Tx
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
	maturityHeight int64
	keyIndex       uint32
	isLocked       bool
}

// isMature returns true if the target utxo is considered "mature" at the
// passed block height. Otherwise, false is returned.
func (u *utxo) isMature(height int64) bool {
	return height >= u.maturityHeight
}

// keyToAddr maps the passed private to corresponding p2pkh address.
func keyToAddr(key *secp256k1.PrivateKey, net *chaincfg.Params) (dcrutil.Address, error) {
	pubKey := (*secp256k1.PublicKey)(&key.PublicKey)
	serializedKey := pubKey.SerializeCompressed()
	pubKeyAddr, err := dcrutil.NewAddressSecpPubKey(serializedKey, net)
	if err != nil {
		return nil, err
	}
	return pubKeyAddr.AddressPubKeyHash(), nil
}
