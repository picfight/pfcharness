package pfcharness

import (
	"github.com/picfight/pfcd/chaincfg"
	"github.com/picfight/pfcd/pfcutil"
	"github.com/jfixby/coinharness"
)

type Address struct {
	Address pfcutil.Address
}

func (c *Address) String() string {
	return c.Address.String()
}

func (c *Address) Internal() interface{} {
	return c.Address
}

func (c *Address) IsForNet(net coinharness.Network) bool {
	return c.Address.IsForNet(net.Params().(*chaincfg.Params))
}
