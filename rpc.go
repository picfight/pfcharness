package pfcharness

import (
	"fmt"
	"github.com/jfixby/coinharness"
	"github.com/jfixby/pin"
	"github.com/picfight/pfcd/chaincfg"
	"github.com/picfight/pfcd/pfcjson"
	"github.com/picfight/pfcd/pfcutil"
	"github.com/picfight/pfcd/rpcclient"
	"io/ioutil"
)

type PfcRPCClientFactory struct {
}

func (f *PfcRPCClientFactory) NewRPCConnection(config coinharness.RPCConnectionConfig, handlers coinharness.RPCClientNotificationHandlers) (coinharness.RPCClient, error) {
	var h *rpcclient.NotificationHandlers
	if handlers != nil {
		h = handlers.(*rpcclient.NotificationHandlers)
	}

	file := config.CertificateFile
	fmt.Println("reading: " + file)
	cert, err := ioutil.ReadFile(file)
	pin.CheckTestSetupMalfunction(err)

	cfg := &rpcclient.ConnConfig{
		Host:                 config.Host,
		Endpoint:             config.Endpoint,
		User:                 config.User,
		Pass:                 config.Pass,
		Certificates:         cert,
		DisableAutoReconnect: true,
		HTTPPostMode:         false,
	}

	return NewRPCClient(cfg, h)
}

type PFCPCClient struct {
	rpc *rpcclient.Client
}

func (c *PFCPCClient) AddNode(args *coinharness.AddNodeArguments) error {
	return c.rpc.AddNode(args.TargetAddr, args.Command.(rpcclient.AddNodeCommand))
}

func (c *PFCPCClient) Disconnect() {
	c.rpc.Disconnect()
}

func (c *PFCPCClient) Shutdown() {
	c.rpc.Shutdown()
}

func (c *PFCPCClient) NotifyBlocks() error {
	return c.rpc.NotifyBlocks()
}

func (c *PFCPCClient) GetBlockCount() (int64, error) {
	return c.rpc.GetBlockCount()
}

func (c *PFCPCClient) Generate(blocks uint32) (result []coinharness.Hash, e error) {
	list, e := c.rpc.Generate(blocks)
	if e != nil {
		return nil, e
	}
	for _, el := range list {
		result = append(result, el)
	}
	return result, nil
}

func (c *PFCPCClient) Internal() interface{} {
	return c.rpc
}

func (c *PFCPCClient) GetRawMempool(command interface{}) (result []coinharness.Hash, e error) {
	list, e := c.rpc.GetRawMempool(command.(pfcjson.GetRawMempoolTxTypeCmd))
	if e != nil {
		return nil, e
	}
	for _, el := range list {
		result = append(result, el)
	}
	return result, nil
}

func (c *PFCPCClient) SendRawTransaction(tx coinharness.CreatedTransactionTx, allowHighFees bool) (result coinharness.Hash, e error) {
	txx := TransactionTxToRaw(tx)
	r, e := c.rpc.SendRawTransaction(txx, allowHighFees)
	return r, e
}

func (c *PFCPCClient) GetPeerInfo() ([]coinharness.PeerInfo, error) {
	pif, err := c.rpc.GetPeerInfo()
	if err != nil {
		return nil, err
	}

	l := []coinharness.PeerInfo{}
	for _, i := range pif {
		inf := coinharness.PeerInfo{}
		inf.Addr = i.Addr
		l = append(l, inf)

	}
	return l, nil
}

func NewRPCClient(config *rpcclient.ConnConfig, handlers *rpcclient.NotificationHandlers) (coinharness.RPCClient, error) {
	legacy, err := rpcclient.New(config, handlers)
	if err != nil {
		return nil, err
	}

	result := &PFCPCClient{rpc: legacy}
	return result, nil
}

func (c *PFCPCClient) GetNewAddress(account string) (coinharness.Address, error) {
	legacy, err := c.rpc.GetNewAddress(account)
	if err != nil {
		return nil, err
	}

	result := &PFCAddress{Address: legacy}
	return result, nil
}

func (c *PFCPCClient) ValidateAddress(address coinharness.Address) (*coinharness.ValidateAddressResult, error) {
	legacy, err := c.rpc.ValidateAddress(address.Internal().(pfcutil.Address))
	// *pfcjson.ValidateAddressWalletResult
	if err != nil {
		return nil, err
	}
	result := &coinharness.ValidateAddressResult{
		Address:      legacy.Address,
		Account:      legacy.Account,
		IsValid:      legacy.IsValid,
		IsMine:       legacy.IsMine,
		IsCompressed: legacy.IsCompressed,
	}
	return result, nil
}

func (c *PFCPCClient) GetBalance(account string) (*coinharness.GetBalanceResult, error) {
	legacy, err := c.rpc.GetBalance(account)
	// *pfcjson.ValidateAddressWalletResult
	if err != nil {
		return nil, err
	}
	result := &coinharness.GetBalanceResult{
		BlockHash:      legacy.BlockHash,
		TotalSpendable: legacy.TotalSpendable,
	}
	return result, nil
}

func (c *PFCPCClient) GetBestBlock() (coinharness.Hash, int64, error) {
	return c.rpc.GetBestBlock()
}

func (c *PFCPCClient) CreateNewAccount(account string) error {
	return c.rpc.CreateNewAccount(account)
}

func (c *PFCPCClient) WalletLock() error {
	return c.rpc.WalletLock()
}

func (c *PFCPCClient) WalletInfo() (*coinharness.WalletInfoResult, error) {
	r, err := c.rpc.WalletInfo()
	if err != nil {
		return nil, err
	}
	result := &coinharness.WalletInfoResult{
		Unlocked:        r.Unlocked,
		DaemonConnected: r.DaemonConnected,
		Voting:          r.DaemonConnected,
	}
	return result, nil
}

func (c *PFCPCClient) WalletUnlock(passphrase string, timeoutSecs int64) error {
	return c.rpc.WalletPassphrase(passphrase, timeoutSecs)
}

func (c *PFCPCClient) CreateTransaction(*coinharness.CreateTransactionArgs) (coinharness.CreatedTransactionTx, error) {
	panic("")
}

func (c *PFCPCClient) GetBuildVersion() (coinharness.BuildVersion, error) {
	//legacy, err := c.rpc.GetBuildVersion()
	//if err != nil {
	//	return nil, err
	//}
	//return legacy, nil
	return nil, fmt.Errorf("decred does not support this feature (GetBuildVersion)")
}

type PFCAddress struct {
	Address pfcutil.Address
}

func (c *PFCAddress) String() string {
	return c.Address.String()
}

func (c *PFCAddress) Internal() interface{} {
	return c.Address
}

func (c *PFCAddress) IsForNet(net coinharness.Network) bool {
	return c.Address.IsForNet(net.(*chaincfg.Params))
}
