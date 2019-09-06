package walletcls

import (
	"github.com/jfixby/coinharness"
	"github.com/jfixby/coinharness/consolewallet"
	"github.com/jfixby/pin"
	"github.com/jfixby/pin/commandline"
	"github.com/picfight/pfcharness"
	"path/filepath"
)

// ConsoleWalletFactory produces a new ConsoleWallet-instance upon request
type ConsoleWalletFactory struct {
	// WalletExecutablePathProvider returns path to the btcd executable
	WalletExecutablePathProvider commandline.ExecutablePathProvider
	ConsoleCommandCook           PfcdConsoleCommandCook
	RPCClientFactory             pfcharness.PfcRPCClientFactory
}

// NewWallet creates and returns a fully initialized instance of the ConsoleWallet.
func (factory *ConsoleWalletFactory) NewWallet(config *coinharness.TestWalletConfig) coinharness.Wallet {
	pin.AssertNotNil("ActiveNet", config.ActiveNet)
	pin.AssertNotNil("WorkingDir", config.WorkingDir)
	pin.AssertNotEmpty("WorkingDir", config.WorkingDir)

	args := &consolewallet.NewConsoleWalletArgs{
		ClientFac:                    &factory.RPCClientFactory,
		ConsoleCommandCook:           &factory.ConsoleCommandCook,
		WalletExecutablePathProvider: factory.WalletExecutablePathProvider,
		RpcUser:                      "user",
		RpcPass:                      "pass",
		AppDir:                       filepath.Join(config.WorkingDir, "pfcd"),
		WalletRPCHost:                config.WalletRPCHost,
		WalletRPCPort:                config.WalletRPCPort,
		ActiveNet:                    config.ActiveNet,
	}

	return consolewallet.NewConsoleWallet(args)
}

type PfcdConsoleCommandCook struct {
}

// cookArguments prepares arguments for the command-line call
func (cook *PfcdConsoleCommandCook) CookArguments(par *consolewallet.ConsoleCommandParams) map[string]interface{} {
	result := make(map[string]interface{})

	result["txindex"] = commandline.NoArgumentValue
	result["addrindex"] = commandline.NoArgumentValue
	result["rpcuser"] = par.RpcUser
	result["rpcpass"] = par.RpcPass
	result["rpcconnect"] = par.RpcConnect
	result["rpclisten"] = par.RpcListen
	result["datadir"] = par.AppDir
	result["debuglevel"] = par.DebugLevel
	result["profile"] = par.Profile
	result["rpccert"] = par.CertFile
	result["rpckey"] = par.KeyFile
	result[pfcharness.NetworkFor(par.Network)] = commandline.NoArgumentValue

	commandline.ArgumentsCopyTo(par.ExtraArguments, result)
	return result
}
