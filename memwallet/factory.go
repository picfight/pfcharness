package memwallet

import (
	"github.com/jfixby/coinharness"
	"github.com/jfixby/pin"
	"github.com/picfight/pfcd/chaincfg"
	"github.com/picfight/pfcd/chaincfg/chainhash"
	"github.com/picfight/pfcd/wire"
	"github.com/picfight/pfcharness"
	"github.com/picfight/pfcutil"
	"github.com/picfight/pfcutil/hdkeychain"
)

// MemWalletFactory produces a new InMemoryWallet-instance upon request
type MemWalletFactory struct {
}

// NewWallet creates and returns a fully initialized instance of the InMemoryWallet.
func (f *MemWalletFactory) NewWallet(cfg *coinharness.TestWalletConfig) coinharness.Wallet {
	pin.AssertNotNil("ActiveNet", cfg.ActiveNet)
	w, e := newMemWallet(cfg.ActiveNet.(*chaincfg.Params), cfg.Seed.([chainhash.HashSize + 4]byte))
	pin.CheckTestSetupMalfunction(e)
	return w
}

func newMemWallet(net *chaincfg.Params, harnessHDSeed [chainhash.HashSize + 4]byte) (*InMemoryWallet, error) {
	hdRoot, err := hdkeychain.NewMaster(harnessHDSeed[:], net)
	if err != nil {
		return nil, nil
	}

	// The first child key from the hd root is reserved as the coinbase
	// generation address.
	coinbaseChild, err := hdRoot.Child(0)
	if err != nil {
		return nil, err
	}
	coinbaseKey, err := coinbaseChild.ECPrivKey()
	if err != nil {
		return nil, err
	}
	coinbaseAddr, err := keyToAddr(coinbaseKey, net)
	if err != nil {
		return nil, err
	}

	// Track the coinbase generation address to ensure we properly track
	// newly generated coins we can spend.
	addrs := make(map[uint32]pfcutil.Address)
	addrs[0] = coinbaseAddr

	clientFac := &pfcharness.PfcRPCClientFactory{}

	return &InMemoryWallet{
		net:               net,
		coinbaseKey:       coinbaseKey,
		coinbaseAddr:      coinbaseAddr,
		hdIndex:           1,
		hdRoot:            hdRoot,
		addrs:             addrs,
		utxos:             make(map[wire.OutPoint]*utxo),
		chainUpdateSignal: make(chan string),
		reorgJournal:      make(map[int64]*undoEntry),
		RPCClientFactory:           clientFac,
	}, nil
}
