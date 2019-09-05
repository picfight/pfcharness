// Copyright (c) 2018 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package pfcharness

import (
	"fmt"
	"github.com/jfixby/coinharness"
	"github.com/jfixby/pin"
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

type PFCRPCClient struct {
	rpc *rpcclient.Client
}

func (c *PFCRPCClient) AddNode(args *coinharness.AddNodeArguments) error {
	return c.rpc.AddNode(args.TargetAddr, args.Command.(rpcclient.AddNodeCommand))
}

func (c *PFCRPCClient) Disconnect() () {
	c.rpc.Disconnect()
}

func (c *PFCRPCClient) Shutdown() () {
	c.rpc.Shutdown()
}

func (c *PFCRPCClient) NotifyBlocks() (error) {
	return c.rpc.NotifyBlocks()
}

func (c *PFCRPCClient) GetBlockCount() (int64, error) {
	return c.rpc.GetBlockCount()
}

func (c *PFCRPCClient) Generate(blocks uint32) (result []coinharness.Hash, e error) {
	list, e := c.rpc.Generate(blocks)
	if e != nil {
		return nil, e
	}
	for _, el := range list {
		result = append(result, el)
	}
	return result, nil
}

func (c *PFCRPCClient) Internal() (interface{}) {
	return c.rpc
}

func (c *PFCRPCClient) GetRawMempool() (result []coinharness.Hash, e error) {
	list, e := c.rpc.GetRawMempool()
	if e != nil {
		return nil, e
	}
	for _, el := range list {
		result = append(result, el)
	}
	return result, nil
}

func (c *PFCRPCClient) SendRawTransaction(tx coinharness.CreatedTransactionTx, allowHighFees bool) (result coinharness.Hash, e error) {
	txx := TransactionTxToRaw(tx)
	r, e := c.rpc.SendRawTransaction(txx, allowHighFees)
	return r, e
}

func (c *PFCRPCClient) GetPeerInfo() ([]coinharness.PeerInfo, error) {
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

	result := &PFCRPCClient{rpc: legacy}
	return result, nil
}
