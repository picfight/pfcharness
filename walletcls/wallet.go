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
	ConsoleCommandCook           WalletConsoleCommandCook
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
		AppDir:                       filepath.Join(config.WorkingDir, "pfcwallet"),
		NodeRPCHost:                  config.NodeRPCHost,
		NodeRPCPort:                  config.NodeRPCPort,
		ActiveNet:                    config.ActiveNet,
	}

	return consolewallet.NewConsoleWallet(args)
}

type WalletConsoleCommandCook struct {
}

// cookArguments prepares arguments for the command-line call
func (cook *WalletConsoleCommandCook) CookArguments(par *consolewallet.ConsoleCommandParams) map[string]interface{} {
	result := make(map[string]interface{})

	result["username"] = par.RpcUser
	result["password"] = par.RpcPass
	result["rpcconnect"] = par.RpcConnect
	result["rpclisten"] = par.RpcListen
	result["appdata"] = par.AppDir
	result["debuglevel"] = par.DebugLevel
	result["cafile"] = par.NodeCertFile
	result["rpccert"] = par.CertFile
	result["rpckey"] = par.KeyFile

	commandline.ArgumentsCopyTo(par.ExtraArguments, result)
	return result
}
