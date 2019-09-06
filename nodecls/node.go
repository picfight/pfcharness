package nodecls

import (
	"github.com/jfixby/coinharness"
	"github.com/jfixby/coinharness/consolenode"
	"github.com/jfixby/pin"
	"github.com/jfixby/pin/commandline"
	"github.com/picfight/pfcharness"
	"path/filepath"
)

// ConsoleNodeFactory produces a new ConsoleNode-instance upon request
type ConsoleNodeFactory struct {
	// NodeExecutablePathProvider returns path to the btcd executable
	NodeExecutablePathProvider commandline.ExecutablePathProvider
	ConsoleCommandCook         PfcdConsoleCommandCook
	RPCClientFactory           pfcharness.PfcRPCClientFactory
}

// NewNode creates and returns a fully initialized instance of the ConsoleNode.
func (factory *ConsoleNodeFactory) NewNode(config *coinharness.TestNodeConfig) coinharness.Node {
	pin.AssertNotNil("WorkingDir", config.WorkingDir)
	pin.AssertNotEmpty("WorkingDir", config.WorkingDir)

	args := &consolenode.NewConsoleNodeArgs{
		ClientFac:                  &factory.RPCClientFactory,
		ConsoleCommandCook:         &factory.ConsoleCommandCook,
		NodeExecutablePathProvider: factory.NodeExecutablePathProvider,
		RpcUser:                    "user",
		RpcPass:                    "pass",
		AppDir:                     filepath.Join(config.WorkingDir, "pfcd"),
		P2PHost:                    config.P2PHost,
		P2PPort:                    config.P2PPort,
		NodeRPCHost:                config.NodeRPCHost,
		NodeRPCPort:                config.NodeRPCPort,
		ActiveNet:                  config.ActiveNet,
	}

	return consolenode.NewConsoleNode(args)
}

type PfcdConsoleCommandCook struct {
}

// cookArguments prepares arguments for the command-line call
func (cook *PfcdConsoleCommandCook) CookArguments(par *consolenode.ConsoleCommandParams) map[string]interface{} {
	result := make(map[string]interface{})

	result["txindex"] = commandline.NoArgumentValue
	result["addrindex"] = commandline.NoArgumentValue
	result["rpcuser"] = par.RpcUser
	result["rpcpass"] = par.RpcPass
	result["rpcconnect"] = par.RpcConnect
	result["rpclisten"] = par.RpcListen
	result["listen"] = par.P2pAddress
	result["datadir"] = par.AppDir
	result["debuglevel"] = par.DebugLevel
	result["profile"] = par.Profile
	result["rpccert"] = par.CertFile
	result["rpckey"] = par.KeyFile
	if par.MiningAddress != nil {
		result["miningaddr"] = par.MiningAddress.String()
	}
	result[networkFor(par.Network)] = commandline.NoArgumentValue

	commandline.ArgumentsCopyTo(par.ExtraArguments, result)
	return result
}

// networkFor resolves network argument for node and wallet console commands
func networkFor(net coinharness.Network) string {
	if net == &chaincfg.SimNetParams {
		return "simnet"
	}
	if net == &chaincfg.TestNet3Params {
		return "testnet"
	}
	if net == &chaincfg.RegressionNetParams {
		return "regtest"
	}
	if net == &chaincfg.MainNetParams {
		// no argument needed for the MainNet
		return commandline.NoArgument
	}

	// should never reach this line, report violation
	pin.ReportTestSetupMalfunction(fmt.Errorf("unknown network: %v ", net))
	return ""
}
