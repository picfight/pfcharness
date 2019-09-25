package pfcharness

import (
	"fmt"
	"github.com/jfixby/coinharness"
	"github.com/jfixby/pin"
	"github.com/jfixby/pin/commandline"
	"github.com/picfight/pfcd/chaincfg"
)

type Network struct {
	Net *chaincfg.Params
}

func (n *Network) Params() interface{} {
	return n.Net
}

func (n *Network) CoinbaseMaturity() int64 {
	return int64(n.Net.CoinbaseMaturity)
}

// networkFor resolves network argument for node and wallet console commands
func NetworkFor(net coinharness.Network) string {
	if net.Params() == &chaincfg.SimNetParams {
		return "simnet"
	}
	if net.Params() == &chaincfg.TestNet3Params {
		return "testnet"
	}
	if net.Params() == &chaincfg.RegNetParams {
		return "regnet"
	}
	if net.Params() == &chaincfg.PicFightCoinNetParams {
		return commandline.NoArgument
	}
	if net.Params() == &chaincfg.DecredNetParams {
		panic("network is not supported")
	}

	// should never reach this line, report violation
	pin.ReportTestSetupMalfunction(fmt.Errorf("unknown network: %v ", net))
	return ""
}
