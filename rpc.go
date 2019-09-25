package pfcharness

import (
	"fmt"
	"github.com/picfight/pfcd/chaincfg/chainhash"
	"github.com/picfight/pfcd/pfcjson"
	"github.com/picfight/pfcd/pfcutil"
	"github.com/picfight/pfcd/rpcclient"
	"github.com/jfixby/coinharness"
	"github.com/jfixby/pin"
	"io/ioutil"
)

type RPCClientFactory struct {
}

func (f *RPCClientFactory) NewRPCConnection(config coinharness.RPCConnectionConfig, handlers coinharness.RPCClientNotificationHandlers) (coinharness.RPCClient, error) {
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

func NewRPCClient(config *rpcclient.ConnConfig, handlers *rpcclient.NotificationHandlers) (coinharness.RPCClient, error) {
	legacy, err := rpcclient.New(config, handlers)
	if err != nil {
		return nil, err
	}

	result := &RPCClient{rpc: legacy}
	return result, nil
}

type RPCClient struct {
	rpc *rpcclient.Client
}

func (c *RPCClient) ListUnspent() ([]*coinharness.Unspent, error) {
	result, err := c.rpc.ListUnspent()
	if err != nil {
		return nil, err
	}
	var r []*coinharness.Unspent
	for _, e := range result {
		x := &coinharness.Unspent{}

		x.TxID = e.TxID
		x.Vout = e.Vout
		x.Tree = e.Tree
		x.TxType = e.TxType
		x.Address = e.Address
		x.Account = e.Account
		x.ScriptPubKey = e.ScriptPubKey
		x.RedeemScript = e.RedeemScript
		x.Amount = coinharness.CoinsAmountFromFloat(e.Amount)
		x.Confirmations = e.Confirmations
		x.Spendable = e.Spendable

		r = append(r, x)
	}

	return r, nil
}

func (c *RPCClient) AddNode(args *coinharness.AddNodeArguments) error {
	return c.rpc.AddNode(args.TargetAddr, args.Command.(rpcclient.AddNodeCommand))
}

func (c *RPCClient) LoadTxFilter(reload bool, addr []coinharness.Address) error {
	addresses := []pfcutil.Address{}
	for _, e := range addr {
		addresses = append(addresses, e.Internal().(pfcutil.Address))
	}
	return c.rpc.LoadTxFilter(reload, addresses, nil)
}

func (c *RPCClient) SubmitBlock(block coinharness.Block) error {
	return c.rpc.SubmitBlock(block.(*pfcutil.Block), nil)
}

func (c *RPCClient) Disconnect() {
	c.rpc.Disconnect()
}

func (c *RPCClient) Shutdown() {
	c.rpc.Shutdown()
}

func (c *RPCClient) NotifyBlocks() error {
	return c.rpc.NotifyBlocks()
}

func (c *RPCClient) GetBlockCount() (int64, error) {
	return c.rpc.GetBlockCount()
}

func (c *RPCClient) Generate(blocks uint32) (result []coinharness.Hash, e error) {
	list, e := c.rpc.Generate(blocks)
	if e != nil {
		return nil, e
	}
	for _, el := range list {
		result = append(result, el)
	}
	return result, nil
}

func (c *RPCClient) Internal() interface{} {
	return c.rpc
}

func (c *RPCClient) GetRawMempool(command interface{}) (result []coinharness.Hash, e error) {
	list, e := c.rpc.GetRawMempool(command.(pfcjson.GetRawMempoolTxTypeCmd))
	if e != nil {
		return nil, e
	}
	for _, el := range list {
		result = append(result, el)
	}
	return result, nil
}

func (c *RPCClient) SendRawTransaction(tx *coinharness.MessageTx, allowHighFees bool) (result coinharness.Hash, e error) {
	txx := TransactionTxToRaw(tx)
	r, e := c.rpc.SendRawTransaction(txx, allowHighFees)
	return r, e
}

func (c *RPCClient) GetBlock(hash coinharness.Hash) (*coinharness.MsgBlock, error) {
	block, err := c.rpc.GetBlock(hash.(*chainhash.Hash)) //*wire.MsgBlock
	if err != nil {
		return nil, err
	}

	b := &coinharness.MsgBlock{}

	ttx := []*coinharness.MessageTx{}
	for _, ti := range block.Transactions {
		ttx = append(ttx,
			TransactionRawToTx(ti),
		)
	}

	b.Transactions = ttx

	return b, nil
}

func (c *RPCClient) GetPeerInfo() ([]coinharness.PeerInfo, error) {
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

func (c *RPCClient) GetNewAddress(account string) (coinharness.Address, error) {
	legacy, err := c.rpc.GetNewAddress(account)
	if err != nil {
		return nil, err
	}

	result := &Address{Address: legacy}
	return result, nil
}

func (c *RPCClient) ValidateAddress(address coinharness.Address) (*coinharness.ValidateAddressResult, error) {
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

func (c *RPCClient) GetBalance() (*coinharness.GetBalanceResult, error) {
	legacy, err := c.rpc.GetBalance("*")
	// *pfcjson.ValidateAddressWalletResult
	if err != nil {
		return nil, err
	}
	result := &coinharness.GetBalanceResult{
		BlockHash: legacy.BlockHash,
		//TotalSpendable:   coinharness.CoinsAmountFromFloat(legacy.TotalSpendable),
		//TotalUnconfirmed: coinharness.CoinsAmountFromFloat(legacy.TotalUnconfirmed),
		//
		//CumulativeTotal:              coinharness.CoinsAmountFromFloat(legacy.CumulativeTotal),
		//TotalVotingAuthority:         coinharness.CoinsAmountFromFloat(legacy.TotalVotingAuthority),
		//TotalLockedByTickets:         coinharness.CoinsAmountFromFloat(legacy.TotalLockedByTickets),
		//TotalImmatureStakeGeneration: coinharness.CoinsAmountFromFloat(legacy.TotalImmatureStakeGeneration),
		//TotalImmatureCoinbaseRewards: coinharness.CoinsAmountFromFloat(legacy.TotalImmatureCoinbaseRewards),
	}
	result.Balances = make(map[string]coinharness.GetAccountBalanceResult)
	for _, v := range legacy.Balances {
		x := coinharness.GetAccountBalanceResult{
			AccountName:             v.AccountName,
			Total:                   coinharness.CoinsAmountFromFloat(v.Total),
			Spendable:               coinharness.CoinsAmountFromFloat(v.Spendable),
			Unconfirmed:             coinharness.CoinsAmountFromFloat(v.Unconfirmed),
			LockedByTickets:         coinharness.CoinsAmountFromFloat(v.LockedByTickets),
			VotingAuthority:         coinharness.CoinsAmountFromFloat(v.VotingAuthority),
			ImmatureCoinbaseRewards: coinharness.CoinsAmountFromFloat(v.ImmatureCoinbaseRewards),
		}
		result.Balances[v.AccountName] = x
	}

	return result, nil
}

func (c *RPCClient) GetBestBlock() (coinharness.Hash, int64, error) {
	return c.rpc.GetBestBlock()
}

func (c *RPCClient) ListAccounts() (map[string]coinharness.CoinsAmount, error) {
	l, err := c.rpc.ListAccounts()
	if err != nil {
		return nil, err
	}

	r := make(map[string]coinharness.CoinsAmount)
	for k, v := range l {
		r[k] = coinharness.CoinsAmount{int64(v)}
	}
	return r, nil
}

func (c *RPCClient) CreateNewAccount(account string) error {
	return c.rpc.CreateNewAccount(account)
}

func (c *RPCClient) WalletLock() error {
	return c.rpc.WalletLock()
}

func (c *RPCClient) WalletInfo() (*coinharness.WalletInfoResult, error) {
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

func (c *RPCClient) WalletUnlock(passphrase string, timeoutSecs int64) error {
	return c.rpc.WalletPassphrase(passphrase, timeoutSecs)
}

func (c *RPCClient) GetBuildVersion() (coinharness.BuildVersion, error) {
	//legacy, err := c.rpc.GetBuildVersion()
	//if err != nil {
	//	return nil, err
	//}
	//return legacy, nil
	return nil, fmt.Errorf("decred does not support this feature (GetBuildVersion)")
}
