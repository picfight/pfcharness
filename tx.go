package pfcharness

import (
	"github.com/jfixby/coinharness"
	"github.com/picfight/pfcd/wire"
)

type OutputTx struct {
	Parent *wire.TxOut
}

func (o *OutputTx) PkScript() []byte {
	return o.Parent.PkScript
}

func (o *OutputTx) Value() int64 {
	return o.Parent.Value
}

type InputTx struct {
	Parent *wire.TxIn
}

type CreatedTransactionTx struct {
	version  int32
	txIn     []*InputTx
	txOut    []*OutputTx
	lockTime uint32
	txHash   coinharness.Hash
}

func (o *CreatedTransactionTx) LockTime() uint32 {
	return o.lockTime
}

func (o *CreatedTransactionTx) Version() int32 {
	return o.version
}

func (o *CreatedTransactionTx) TxHash() coinharness.Hash {
	return o.txHash
}

func (o *CreatedTransactionTx) TxIn() (result []coinharness.InputTx) {
	result = []coinharness.InputTx{}
	for _, ti := range o.txIn {
		result = append(result, ti)
	}
	return
}

func (o *CreatedTransactionTx) TxOut() (result []coinharness.OutputTx) {
	result = []coinharness.OutputTx{}
	for _, ti := range o.txOut {
		result = append(result, ti)
	}
	return
}
