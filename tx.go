package pfcharness

import (
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
